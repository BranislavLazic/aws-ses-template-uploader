package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsDefaultRegion   = os.Getenv("AWS_DEFAULT_REGION")
)

type Template struct {
	TemplateName string `json:"TemplateName"`
	SubjectPart  string `json:"SubjectPart"`
	HTMLPart     string `json:"HtmlPart"`
	TextPart     string `json:"TextPart"`
}

func newTemplate() *Template {
	return &Template{
		TemplateName: "",
		SubjectPart:  "",
		HTMLPart:     "",
		TextPart:     "",
	}
}

func parseArgs(awsSES *ses.SES, args []string) error {

	switch args[1] {
	case "list":
		return handleListTemplates(awsSES)

	case "create":
		template := newTemplate()
		if len(args) > 2 {
			file, err := os.Open(args[2])
			if err != nil {
				log.Fatalf("file does not exist")
				os.Exit(1)
			}
			byteValue, _ := ioutil.ReadAll(file)
			json.Unmarshal(byteValue, &template)
			return handleCreateTemplate(awsSES, &ses.Template{
				TemplateName: &template.TemplateName,
				SubjectPart:  &template.SubjectPart,
				HtmlPart:     &template.HTMLPart,
				TextPart:     &template.TextPart,
			})
		}
		return errors.New("provide a template as JSON file")
	case "delete":
		if len(args) > 2 {
			return handleDeleteTemplate(awsSES, args[2])
		}
		return errors.New("provide a template name")
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
	for _, data := range out.TemplatesMetadata {
		fmt.Printf("Name: %s, Created: %v\n", *data.Name, data.CreatedTimestamp)
	}
	return nil
}

func handleCreateTemplate(awsSES *ses.SES, template *ses.Template) error {
	_, err := awsSES.CreateTemplate(&ses.CreateTemplateInput{
		Template: template,
	})
	if err != nil {
		return err
	}
	log.Printf("successfully created %s template", *template.TemplateName)
	return nil
}
func handleDeleteTemplate(awsSES *ses.SES, templateName string) error {
	_, err := awsSES.DeleteTemplate(&ses.DeleteTemplateInput{TemplateName: &templateName})
	if err != nil {
		return err
	}
	log.Printf("successfully deleted %s template", templateName)
	return nil
}

func main() {
	// Command arguments: list, create, delete... with template name or template file.
	cmdArgs := os.Args
	// Open AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(awsDefaultRegion)},
	)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Parse command line arguments and provide AWS SES value
	err = parseArgs(ses.New(awsSession), cmdArgs)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
