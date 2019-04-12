@setlocal enableextensions
@cd /d "%~dp0"
echo New Relic OHI Installer

@echo off
goto do_Install

:do_Install
    net session >nul 2>&1
    if %errorLevel% == 0 (
        echo Success: Administrative permissions confirmed.
        copy nri-flex.exe "C:\Program Files\New Relic\newrelic-infra\custom-integrations\"
        copy nri-flex-def-win.yml "C:\Program Files\New Relic\newrelic-infra\custom-integrations\"
        copy nri-flex-config.yml "C:\Program Files\New Relic\newrelic-infra\integrations.d\"

        mkdir "C:\Program Files\New Relic\newrelic-infra\custom-integrations\flexConfigs"
        mkdir "C:\Program Files\New Relic\newrelic-infra\custom-integrations\flexContainerDiscovery"

        REM dont forget to copy your configs into the relevant directories and nrjmx if wanting to use jmx

        net stop newrelic-infra
        net start newrelic-infra
    ) else (
        echo Failure: Administrative permissions required!
    )

timeout 5