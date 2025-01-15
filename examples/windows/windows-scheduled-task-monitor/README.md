This integration monitors the windows scheduled tasks execution status and missed executions.

This integration contains two parts,
1. a flex based integration to gather execution status and missed executions and report as a custom event.
2. captures scheduled job history as windows event logs and report them to NR Log management solution.

Follow the below steps to configure this integration to monitor windows scheduled tasks:

Pre-requisites:
1. A new relic account
2. New relic infrastructure agent set up on the host
3. Power shell installed on the host
4. Task Scheduler history logs are enabled on the host (https://docs.nxlog.co/integrate/windows-task-scheduler.html)

Configure the integration
1. Copy below files to C:\Program Files\New Relic\newrelic-infra\integrations.d\ directory
    a. WindowsScheduledTaskInfo.ps1
    b. windows-scheduledtask-info-new.yml
2. Copy below files to C:\Program Files\New Relic\newrelic-infra\logging.d\ directory
    a. winevt_taskscheduler_logs.yml

Configure Dashboard to view the data
1. Import below csv files as lookup tables
    a. Refer to https://docs.newrelic.com/docs/logs/ui-data/lookup-tables-ui/
    b. windows-event-log-level.csv with Table name as windows_event_log_level
    c. windows-task-result-descriptions.csv with Table name as windows_scheduled_task_result_descriptions
2. Import Dashboard
    a. open dashboard_windows_scheduled_task_oveview.json file in a text editor, and replace accountId 9999999 with the accountId, where you are importing this dashboard.
    b. On New Relic UI, Navigate to Dashboards and click on “Import Dashboard”.
    c. Copy the updated json content from the above step in the “Paste your json code” text box and click on the Import dashboard button.
    d. Edit the dashboard created and update the default value for the Scheduled Task variable.
