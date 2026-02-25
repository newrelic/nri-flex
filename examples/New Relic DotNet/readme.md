# New Relic Flex .NET Applications Plugin Tracker

This New Relic Infrastructure integration uses the Flex integration to collect information about .NET applications in your New Relic account and their installed plugins.

## Overview

This integration performs two main tasks:
1. Retrieves a list of all .NET applications reporting to New Relic
2. Extracts plugin information for each .NET application

The collected data is sent to New Relic where it can be queried and visualized.

## Requirements

- New Relic Infrastructure Agent installed
- New Relic Flex integration installed
- A valid New Relic API key with query permissions

## Installation

1. Install the New Relic Infrastructure agent if not already installed
   ```
   # Follow instructions at https://docs.newrelic.com/install/infrastructure/
   ```

2. Install the New Relic Flex integration if not already installed
   ```
   # On Linux
   sudo apt-get install nri-flex

   # Or using the New Relic CLI
   newrelic install -n infrastructure-agent-flex
   ```

3. Save the configuration as a YAML file in the Flex integration directory:
   ```
   # Linux
   sudo mkdir -p /etc/newrelic-infra/integrations.d/flex
   sudo vi /etc/newrelic-infra/integrations.d/flex/dotnet-app-plugins.yml
   
   # Windows
   mkdir "C:\Program Files\New Relic\newrelic-infra\integrations.d\flex"
   notepad "C:\Program Files\New Relic\newrelic-infra\integrations.d\flex\dotnet-app-plugins.yml"
   ```

4. Paste the configuration into the file, making sure to replace `APIKEYHERE` with your actual New Relic API key.

5. Restart the New Relic Infrastructure agent
   ```
   # Linux
   sudo systemctl restart newrelic-infra
   
   # Windows
   Restart-Service -Name "newrelic-infra"
   ```

## Configuration

The configuration consists of two API calls:

### 1. dotnetAppList

This API call queries New Relic's GraphQL API to get a list of all .NET applications that are currently reporting to New Relic.

Key settings:
- Event type: `dotnetAppList`
- GraphQL query retrieves name and GUID for each .NET application
- 5-minute timeout to allow for a large number of applications

### 2. dotnetPluginList

This API call uses the GUIDs obtained from the first call to query detailed information about each .NET application, including:
- Application name
- GUID
- Account ID and name
- Installed plugins
- Host information
- Instance ID
- Any tags applied to the application

The JQ processor transforms the response to create detailed events about each plugin installed on each application instance.

## Data Collection

The integration:
- Runs every 10 minutes (configurable via the `interval` setting)
- Has a 32-second timeout for the integration as a whole
- Has a 320-second timeout for this specific configuration

## Data Usage

Once the data is in New Relic, you can:

1. Query the data in NRQL:
   ```
   FROM dotnetPluginList SELECT appName, guid, library, host 
   ```

2. Create dashboards to visualize:
   - Distribution of plugins across applications
   - Applications with outdated plugins
   - Correlations between plugins and application performance

## Troubleshooting

If the integration is not working:

It's usually a timeout issue... Check that first!

1. Check the Infrastructure agent logs:
   ```
   # Linux
   tail -f /var/log/newrelic-infra/newrelic-infra.log
   
   # Windows
   Get-Content -Path "C:\Program Files\New Relic\newrelic-infra\newrelic-infra.log" -Tail 20
   ```

2. Verify the API key has the necessary permissions