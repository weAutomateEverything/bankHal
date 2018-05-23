package bankSkynet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/chef"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Service interface {
	RecreateNode(ctx context.Context, chatid uint32,nodeName, callerName string) error
	sendSkynetAlert(ctx context.Context,chatid uint32, message string)
}

type service struct {
	alert          alert.Service
	chefStore      chef.Store
	calloutService callout.Service
}

func NewService(alert alert.Service, chefStore chef.Store, calloutService callout.Service) Service {
	return &service{alert, chefStore, calloutService}
}

func (s *service) sendSkynetAlert(ctx context.Context,chatid uint32, message string) {
	if !s.checkSend(ctx, chatid, message) {
		log.Println("Ignoreing message: " + message)
		return
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(message), &dat); err != nil {
		s.alert.SendError(ctx, fmt.Errorf("skynet alert - rrror unmarshalling: %s", message))
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprintf(":computer:"))
	buffer.WriteString(" ")
	buffer.WriteString("*Skynet Alert*\n")
	buffer.WriteString(dat["message"].(string))
	buffer.WriteString("\n")

	args := dat["args"].(map[string]interface{})
	config := args["config"].(map[string]interface{})
	chefConfig := config["chef_config"].(map[string]interface{})
	buffer.WriteString("*Environment: *")
	buffer.WriteString(chefConfig["environment"].(string))
	buffer.WriteString("\n")

	runlist := chefConfig["run_list"].([]interface{})
	log.Println(runlist)
	buffer.WriteString("*Run list: *")

	for _, recipe := range runlist {
		recipeS := strings.Replace(recipe.(string), "recipe[", "", -1)
		recipeS = strings.Replace(recipeS, "]", "", -1)
		buffer.WriteString(recipeS)
		buffer.WriteString(" ")
	}

	buffer.WriteString("\n")

	cmdbConfig := config["cmdb_config"].(map[string]interface{})

	buffer.WriteString("*Description: *")
	buffer.WriteString(cmdbConfig["description"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Technical: *")
	buffer.WriteString(cmdbConfig["technical"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Environment: *")
	buffer.WriteString(cmdbConfig["environment"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Application: *")
	buffer.WriteString(cmdbConfig["application"].(string))
	buffer.WriteString("\n")

	s.alert.SendAlert(ctx, chatid,buffer.String())

}

func (s *service) checkSend(ctx context.Context,chatid uint32, message string) bool {
	message = strings.ToUpper(message)
	log.Printf("Checking if we should send: %s", message)
	recipes, err := s.chefStore.GetRecipes()
	if err != nil {
		s.alert.SendError(ctx, err)
		return false
	}
	for _, recipe := range recipes {
		check := "RECIPE[" + strings.ToUpper(recipe.Recipe) + "]"
		log.Printf("Comparing %s", check)
		if strings.Contains(message, check) {
			log.Printf("Match Found, returning true")
			return true
		}
	}
	log.Printf("No match found, not sending message")
	return false
}

/*
RecreateNode will find a node, get the details, delete the node, then recreate it
*/
func (s service) RecreateNode(ctx context.Context, chatid uint32,nodeName, callerName string) error {

	skynet := getSkynetUrl()
	json, err := s.findNode(ctx,chatid, nodeName, skynet)
	if err != nil {
		return err
	}

	err = s.deleteNode(ctx,chatid, nodeName, callerName, skynet)
	if err != nil {
		return err
	}

	err = s.waitForDelete(ctx,chatid, nodeName, skynet)
	if err != nil {
		return err
	}

	err = s.createNode(ctx,chatid, json, skynet)
	err = s.waitForBuild(ctx,chatid, nodeName, skynet)
	return nil

}

func (s *service) findNode(ctx context.Context, chatid uint32, nodeName string, skynet string) (string, error) {
	body, err := s.doHTTP(ctx,chatid, "GET", skynet+"/virtual_machines/"+nodeName, "", skynet)
	if err != nil {
		return "", err
	}
	return body, nil
}

func (s service) deleteNode(ctx context.Context,chatid uint32, nodeName, callerName string, skynet string) error {
	s.alert.SendAlert(ctx,chatid, fmt.Sprintf("Received a Delete Node request from %s for node %s. Proceeding with Delete", callerName, nodeName))
	_, err := s.doHTTP(ctx,chatid, "DELETE", skynet+"/virtual_machines/"+nodeName, "", skynet)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) waitForDelete(ctx context.Context,chatid uint32, nodeName string, skynet string) error {
	return s.poll(ctx,chatid, "ARCHIVED", nodeName, skynet, true, 300)
}

func (s *service) createNode(ctx context.Context,chatid uint32, json string, skynet string) error {

	body, err := s.doHTTP(ctx,chatid, "POST", skynet+"/virtual_machines", json, skynet)
	if err != nil {
		s.calloutService.InvokeCallout(ctx, chatid,"skynet error creating node", fmt.Sprintf("Json: %s, Error: %s", json, err.Error()), nil)
		return err
	}
	log.Println(body)
	return nil
}

func (s *service) waitForBuild(ctx context.Context,chatid uint32,nodeName string, skynet string) error {
	return s.poll(ctx,chatid, "BOOTSTRAPPED", nodeName, skynet, false, 1200)
}

func (s service) doHTTP(ctx context.Context,chatid uint32, method, url, body string, skynet string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		s.logError(ctx, chatid, fmt.Sprintf("skynet error creating URL to find node. Error %s. Method: %s, URL: %s, Body %s", err.Error(), method, url, body))
		return "", err
	}
	req.SetBasicAuth(getSkynetUser(), getSkynetPassword())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		s.logError(ctx, chatid,fmt.Sprintf("error in skynet call.  Error %s. Method: %s, URL: %s, Body %s", err.Error(), method, url, body))
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	return string(bodyText), nil
}

func (s service) logError(ctx context.Context,chatid uint32, error string) {
	s.alert.SendAlert(ctx,chatid, error)
}

func (s *service) poll(ctx context.Context,chatid uint32, expectedState, nodeName string, skynet string, ignoreFailed bool, timeout int) error {
	i := 0
	for i < timeout {
		body, err := s.doHTTP(ctx,chatid, "GET", skynet+"/virtual_machines/"+nodeName+"/state", "", skynet)
		if err != nil {
			s.logError(ctx, chatid, fmt.Sprintf("skynet error retreiving node state:  Error %s. Node: %s", err.Error(), nodeName))
			return err
		}
		var dat map[string]interface{}
		err = json.Unmarshal([]byte(body), &dat)
		if err != nil {
			s.logError(ctx,chatid, fmt.Sprintf("skynet error unamrhsaling %s to json. Node: %s", body, nodeName))
			return err
		}
		state := dat["state"].(map[string]interface{})["current"].(string)
		if strings.ToUpper(state) == expectedState {
			s.alert.SendAlert(ctx,chatid, fmt.Sprintf("%s has been reached state %s.", nodeName, expectedState))
			return nil
		}
		if !ignoreFailed && strings.ToUpper(state) == "FAILED" {
			s.alert.SendAlert(ctx,chatid, fmt.Sprintf("%s has entered a Failed State.", nodeName))
			s.calloutService.InvokeCallout(ctx,chatid, fmt.Sprintf("Skynet Error rebuilding node %s", nodeName), "Node failed to build successfully", nil)
			return fmt.Errorf("%s has entered a Failed State", nodeName)
		}
		i++
		if i%60 == 0 {
			s.alert.SendError(ctx, fmt.Errorf("still waiting for node %s to reach state %s. Curent state is %s", nodeName, expectedState, state))
		}
		time.Sleep(time.Second)
	}
	s.calloutService.InvokeCallout(ctx,chatid, fmt.Sprintf("Timed out waiting for node %s to enter state %s", nodeName, expectedState), "", nil)
	err := fmt.Errorf("timed out waiting for node %s to %s", nodeName, expectedState)
	s.logError(ctx, chatid,err.Error())
	return err
}

func getSkynetUrl() string {
	return os.Getenv("SKYNET_URL")
}

func getSkynetPassword() string {
	return os.Getenv("SKYNET_PASSWORD")
}
func getSkynetUser() string {
	return os.Getenv("SKYNET_USER")

}
