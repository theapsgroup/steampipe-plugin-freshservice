## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29).
- Recompiled plugin with Go version `1.21`.

## v0.0.4 [2023-05-05]

_Enhancements_

- Recompiled with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05)

## v0.0.3 [2023-02-14]

_Enhancements_

- Recompiled with steampipe-plugin-sdk v5.1.2

_Bug fixes_

- Fixed issue where `department_ids` column on `freshservice_requester` table was erroneously returning null [#42](https://github.com/theapsgroup/steampipe-plugin-freshservice/issues/42)

## v0.0.2 [2022-10-07]

_Enhancements_

- Recompiled with go1.19 and plugin-sdk v4.1.7

## v0.0.1 [2022-08-09]

_What's new?_

- Initial release containing the following tables.
  - [freshservice_agent](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_agent)
  - [freshservice_announcement](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_announcement)
  - [freshservice_asset](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_asset)
  - [freshservice_asset_component](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_asset_component)
  - [freshservice_asset_contract](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_asset_contract)
  - [freshservice_asset_type](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_asset_type)
  - [freshservice_business_hour](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_business_hour)
  - [freshservice_change](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_change)
  - [freshservice_change_note](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_change_note)
  - [freshservice_contract](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_contract)
  - [freshservice_contract_type](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_contract_type)
  - [freshservice_department](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_department)
  - [freshservice_location](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_location)
  - [freshservice_problem](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_problem)
  - [freshservice_problem_note](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_problem_note)
  - [freshservice_problem_task](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_problem_task)
  - [freshservice_problem_timeentry](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_problem_timeentry)
  - [freshservice_product](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_product)
  - [freshservice_purchase_order](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_purchase_order)
  - [freshservice_release](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_release)
  - [freshservice_release_note](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_release_note)
  - [freshservice_release_task](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_release_task)
  - [freshservice_release_timeentry](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_release_timeentry)
  - [freshservice_requester](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_requester)
  - [freshservice_service](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_service)
  - [freshservice_sla_policy](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_sla_policy)
  - [freshservice_software](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_software)
  - [freshservice_software_installation](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_software_installation)
  - [freshservice_software_user](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_software_user)
  - [freshservice_solution_article](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_solution_article)
  - [freshservice_solution_category](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_solution_category)
  - [freshservice_solution_folder](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_solution_folder)
  - [freshservice_ticket](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_ticket)
  - [freshservice_ticket_conversation](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_ticket_conversation)
  - [freshservice_ticket_task](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_ticket_task)
  - [freshservice_ticket_timeentry](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_ticket_timeentry)
  - [freshservice_vendor](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables/freshservice_vendor)
