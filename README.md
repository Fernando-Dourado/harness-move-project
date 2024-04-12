# Harness Move

An utility tool to copy/clone/move a project.

## Requirements

- The tool does not create the org or project target.
- As safety operation the tool do not delete the entities from the source project.
- The `api-key` need to have access to read from the source project and write to the target project.

## Usage

Execute the operation running at the following command in your terminal

```sh
./harness-move-project \
  --api-token <SAT_OR_PAT> \
  --account <account_identifier> \
  --source-org <org_identifier> --source-project <project_identifier> \
  --target-org <org_identifier> --target-project <project_identifier>
```

If the source and target projects has the same identifier, you can suppress the `--target-project` argument.

When the tool try to create and entity on target project that the same identifier already exist, it just ignore the error and keep the execution. Using that strategy you can run it multiple times without side effects.

## Supported Entities

- Variables
- Environments
- Environment Groups
- Infrastructure Definition
- Services
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

## Not Supported Entities

- Secrets
- Triggers
- Connectors
- Service Overrides

## Limitation

- The tool can only fetch 1000 elements of each entity type.

## Contributions

I am to express my gratitude for inspiration to create this tool.

- [Aleksa Arsic](https://github.com/aleksa11010): Thank you for the inspiration! Your creativity is amazing!
- Francisco Junior: I appreciate inspiring me to improve. Your guidance was crucial!
