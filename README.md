<div class="title-block" style="text-align: center;" align="center">

# Harness Move

A utility tool to copy/clone a project.

![](https://img.shields.io/github/v/release/jf781/harness-move-project)
![](https://img.shields.io/github/release-date/jf781/harness-move-project)

</div>

## Install

Download the latest version from [releases page](https://github.com/jf781/harness-move-project/releases/latest)

## Requirements

- The tool does not create the target Organization.  It must be pre-existing
- As safety operation, the tool do not delete the entities from the source project.
- The `api-key` need to have access to read from the source project and write to the target project.
- You can run it multiple times, when the same entity already exists in the target project we ignore it and do not report it as an error.

## Usage

Execute the operation running at the following command in your terminal that will read from a CSV file. 


```sh
./harness-move-project \
  --apiToken <SAT_OR_PAT> \
  --accountId <account_identifier> \
  --csvPath ./exampleCsvFile.csv \
  --baseUrl https://app.harness.io
```

When the tool try to create and entity on target project that the same identifier already exist, it silently ignores the error and continues the execution. Using that strategy you can run it multiple times without side effects.

## CSV File

You can run this against a single or multiple projects by providing a CSV file with the following format:

| CSV column Name | Description | Required |
| --------------- | ----------- | -------- |
| `sourceOrg` | The name of source organization | Yes |
| `sourceProject` | The name of source project | Yes |
| `targetOrg` | The name of target organization | Yes |
| `targetProject` | The name of target project | No |

If the `targetProject` is not provided, the tool will use the `sourceProject` as the target project name.

There is an example CSV file named `exampleCsvFile.csv` that you can update with your project details.

## Supported Entities

- Variables
- Environments
- Environment Groups
- Infrastructure Definition
- Services
- Service Overrides V1
- Templates
- Pipelines
- Input Sets
- Roles
- User Groups
- Service Accounts
- Role Assignments
- Resource Groups
- Connectors
- Triggers & Webhooks (Work in progress)
- File Store (Working in progress)
- Feature Flags
- Feature Flag Targets & Target Groups
- File Store

## Not Supported Entities

- Secrets
- Triggers
- Connectors
- Service Overrides V2

## Future items
- Mark source project as read-only
- Create new project SDK key and save as a secret when project is created

## Limitation

- The tool can only fetch 1000 elements of each entity type.
- Tags are not supported and cannot be copied from the source entity to the target one.

## Contributions

I am to express my gratitude for inspiration to create this tool.

- [Aleksa Arsic](https://github.com/aleksa11010): Thank you for the inspiration! Your creativity is amazing!
- Francisco Junior: I appreciate inspiring me to improve. Your guidance was crucial!
