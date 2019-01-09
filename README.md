# Spotlight Grabber

This is a simple tool for saving copies of the Microsoft Spotlight background
images to the current folder or a folder passed as the first command-line
argument. It only copies files which are JPG or PNG format, with a height
greater than 400 pixels and a width greater than the height.

## Building

Go 1.11 is required. To build:

```powershell
go build -o spotlight-grabber.exe
```

## Creating a scheduled task

It's best used as a task scheduled to run every day since the Spotlight assets
folder is merely a cache and old images will be deleted. This can be easily
configured with the following PowerShell executed with admin rights:

```powershell
$action = New-ScheduledTaskAction -Execute "absolute\path\to\spotlight-grabber.exe" -Argument "absolute\path\to\saved\images\folder"
$trigger = NewScheduledTaskTrigger -Daily -At 9am
Register-ScheduledTask -Action $action -Trigger $trigger -TaskName "Spotlight Grabber" -Description "Daily copying of images from the Spotlight assets folder"
```
