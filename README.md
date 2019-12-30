# AWS SES template uploader

CLI tool for upload of templates for AWS SES

## Get it

Download the [latest release](https://github.com/BranislavLazic/aws-ses-template-uploader/releases) or
checkout project and build it locally by executing `make` command.

## Use it

Export AWS access key id, secret access key and region as env. variables:

`export AWS_ACCESS_KEY_ID=your-access-key-id`

`export AWS_SECRET_ACCESS_KEY=your-secret-access-key`

`export AWS_DEFAULT_REGION=your-default-region`

### Commands

To see list of your templates:

`$ aws-ses-template-uploader list`

To create a template from JSON file:

`$ aws-ses-template-uploader create /path/template.json`

To delete a template:

`$ aws-ses-template-uploader delete template-name`

## Example of template

```
{
  "TemplateName": "test-email-template",
  "SubjectPart": "Test subject",
  "HtmlPart": "<div>Test</div>",
  "TextPart": "Test"
}
```
