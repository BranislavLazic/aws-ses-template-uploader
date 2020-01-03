package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/olekukonko/tablewriter"
)

var (
	awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsDefaultRegion   = os.Getenv("AWS_DEFAULT_REGION")
)

type template struct {
	TemplateName string `json:"TemplateName"`
	SubjectPart  string `json:"SubjectPart"`
	HTMLPart     string `json:"HtmlPart"`
	TextPart     string `json:"TextPart"`
}

func parseCLIArgs(awsSES *ses.SES, args []string) error {
	if len(args) < 2 {
		return errors.New(`provide one of the following command line arguments: 
	list
	create file.json
	update file.json
	delete`)
	}
	switch args[1] {
	case "list":
		return handleListTemplates(awsSES)
	case "create":
		return handleSaveTemplate(awsSES, args, func(awsSES *ses.SES, template *ses.Template) error {
			return handleCreateTemplate(awsSES, template)
		})
	case "update":
		return handleSaveTemplate(awsSES, args, func(awsSES *ses.SES, template *ses.Template) error {
			return handleUpdateTemplate(awsSES, template)
		})
	case "delete":
		return handleDeleteTemplate(awsSES, args)
	default:
		break
	}
	return nil
}

func handleListTemplates(awsSES *ses.SES) error {
	out, err := awsSES.ListTemplates(&ses.ListTemplatesInput{})
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Created"})
	for _, data := range out.TemplatesMetadata {
		table.Append([]string{
			*data.Name,
			data.CreatedTimestamp.Local().Format("2 Jan 2006 15:04:05"),
		})
	}
	table.Render()
	return nil
}

func parseTemplate(jsonFile string) (*ses.Template, error) {
	var template *template
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	fileName := file.Name()
	if !strings.HasSuffix(fileName, ".json") {
		return nil, fmt.Errorf("%s is not a JSON file", fileName)
	}
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &template)
	return &ses.Template{
		TemplateName: &template.TemplateName,
		SubjectPart:  &template.SubjectPart,
		HtmlPart:     &template.HTMLPart,
		TextPart:     &template.TextPart,
	}, nil
}

func handleSaveTemplate(awsSES *ses.SES, args []string, saveFn func(*ses.SES, *ses.Template) error) error {
	if len(args) > 2 {
		template, err := parseTemplate(args[2])
		if err != nil {
			return err
		}
		return saveFn(awsSES, template)
	}
	return errors.New("provide a template as JSON file")
}

func handleCreateTemplate(awsSES *ses.SES, template *ses.Template) error {
	_, err := awsSES.CreateTemplate(&ses.CreateTemplateInput{
		Template: template,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created %s template\n", *template.TemplateName)
	return nil
}

func handleUpdateTemplate(awsSES *ses.SES, template *ses.Template) error {
	_, err := awsSES.UpdateTemplate(&ses.UpdateTemplateInput{
		Template: template,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully updated %s template\n", *template.TemplateName)
	return nil
}

func handleDeleteTemplate(awsSES *ses.SES, args []string) error {
	if len(args) > 2 {
		templateName := args[2]
		_, err := awsSES.DeleteTemplate(&ses.DeleteTemplateInput{
			TemplateName: &templateName,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Successfully deleted %s template\n", templateName)
		return nil
	}
	return errors.New("provide a template name")
}

func main() {
	// CLI arguments: list, create, delete... with template name or template file.
	args := os.Args
	// Open AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(awsDefaultRegion)},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Parse command line arguments and provide AWS SES value
	err = parseCLIArgs(ses.New(awsSession), args)
	if err != nil {
		log.Fatalln(err)
	}
}
