<div class="title-block" style="text-align: center;" align="center">

# Harness Move

An utility tool to copy/clone/move a project.

![](https://img.shields.io/github/v/release/Fernando-Dourado/harness-move-project)
![](https://img.shields.io/github/release-date/Fernando-Dourado/harness-move-project)

</div>

## Install

Download the latest version from [releases page](https://github.com/Fernando-Dourado/harness-move-project/releases/latest)

## Requirements

- The tool does not create the org, can create the project when flag `create-project` is provided.
- As safety operation the tool do not delete the entities from the source project.
- The `api-key` need to have access to read from the source project and write to the target project.
- You can run it multiple times, when the same entity already exists in the target project we ignore it and do not report it as an error.

## Usage

Running that command in your terminal you do the most basic operation.

```bash
./harness-move-project \
  --api-token <SAT_OR_PAT> \
  --account <account_identifier> \
  --source-org <org_identifier> --source-project <project_identifier> \
  --target-org <org_identifier> --target-project <project_identifier>
```

If the source and target projects has the same identifier, you can suppress the `--target-project` argument. Providing `--create-project true` you can create the target project in case it does not exist in the target account and org.

When the tool try to create and entity on target project that the same identifier already exist, it just ignore the error and keep the execution. Using that strategy you can run it multiple times without side effects.

It is also possible to perform the copy between different accounts. To do this, you need to specify the `--target-account` and `--target-token` of the target account.

## Supported Entities

- Variables
- Environments
- Infrastructure Definition
- Services
- Service Overrides V1
- Templates
- Pipelines
- Input Sets
- File Store
- Connectors

## Partial Supported Entities

- Secrets (Text, File, and SSH Credentials)

> **Note:** During the secret copy process, the secret is created in the destination project with a **dummy value**.  
> The actual value is **not** copied. This rule applies to both **secret text** and **secret file** types.  
> You must manually update the secret with the correct value in the destination project after the copy is completed.

## Not Supported Entities

- Secrets (WinRM Credentials)
- Triggers
- Service Overrides V2

## Limitation

- The tool can only fetch 1000 elements of each entity type.
- Tags are not supported and cannot be copied from the source entity to the target one.

## Contributions

I am to express my gratitude for inspiration to create this tool.

- [Aleksa Arsic](https://github.com/aleksa11010): Thank you for the inspiration! Your creativity is amazing!
- Francisco Junior: I appreciate inspiring me to improve. Your guidance was crucial!

## Usage Output

```text
NAME:
   harness-move-project - Non-official Harness CLI to move project between organizations or accounts.

USAGE:
   harness-move-project [options]

VERSION:
   development

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-token value       API authentication token for accessing the source system.
   --account value         The account identifier associated with the source system.
   --source-org value      The organization identifier in the source account.
   --source-project value  The project identifier in the source account.
   --target-org value      The org identifier in the target account.
   --target-project value  The project identifier in the target account.
   --target-token value    API authentication token for accessing the target system.
                           Not needed if target and source accounts are the same.
   --target-account value  The account identifier associated with the target system.
                           Not needed if target and source accounts are the same.
   --create-project value  Creates the project in the target account/org if missing.
   --help, -h              show help
   --version, -v           print the version
```
