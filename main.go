package main

import (
	"fmt"
	"regexp"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/wekex/deps-cleaner/helper"
)

var (
	progressBar           *widget.ProgressBar
	infiniteBar           *widget.ProgressBarInfinite
	progressBarContainer  *fyne.Container
	actionButtonContainer *fyne.Container
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Deps Cleaner")

	content := createContent(myWindow)
	logScroll := createLogScroll()

	split := container.NewVSplit(content, logScroll)

	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

// actionButtons creates a container with two buttons: "Start Cleaning" and "Preview Delete".
// The "Start Cleaning" button triggers the startCleaning function with the current text values of startPathEntry and filterEntry.
// The "Preview Delete" button triggers the listCleaning function with the current text values of startPathEntry and filterEntry.
func actionButtons(startPathEntry, filterEntry *widget.Entry, myWindow fyne.Window) *fyne.Container {
	// Create the "Start Cleaning" button with a click handler that runs startCleaning in a goroutine
	startButton := widget.NewButton("Start Cleaning", func() {
		go startCleaning(startPathEntry.Text, filterEntry.Text, myWindow)
	})

	// Create the "Preview Delete" button with a click handler that runs listCleaning in a goroutine
	viewButton := widget.NewButton("Preview Delete", func() {
		go listCleaning(startPathEntry.Text, filterEntry.Text, myWindow)
	})

	// Return a horizontal box container with the two buttons
	return container.NewHBox(viewButton, startButton)
}

// toggleVisibility function toggles the visibility of the progress bar container and the action button container
// based on the value of the visible parameter.
// If visible is true, the progress bar container is hidden and the action button container is shown.
// If visible is false, the progress bar container is shown and the action button container is hidden.
func toggleVisibility(visible bool) {
	if visible {
		progressBarContainer.Hide()
		actionButtonContainer.Show()
	} else {
		progressBarContainer.Show()
		actionButtonContainer.Hide()
	}
}

// createContent is a function that creates the main content of the application window.
// It takes in a fyne.Window object as a parameter.
func createContent(myWindow fyne.Window) fyne.CanvasObject {

	// Create the base directory label and entry
	baseDirLabel := widget.NewLabel("Base Directory:")
	startPathEntry := widget.NewEntry()
	startPathEntry.SetPlaceHolder("/Enter/Start/Path")

	// Create a button to open the folder dialog for selecting base directory
	startDirButton := widget.NewButton("Select Base Directory", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err == nil && dir != nil {
				startPathEntry.SetText(dir.Path())
			}
		}, myWindow)
	})

	// Create the filter label and entry
	filterLabel := widget.NewLabel("Filter:")
	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Enter Regex filter")

	// Create the progress bars for displaying progress of cleaning
	progressBar = widget.NewProgressBar()
	infiniteBar = widget.NewProgressBarInfinite()

	// Create containers for action buttons and progress bars
	progressBarContainer = container.NewVBox(progressBar, infiniteBar)
	actionButtonContainer = actionButtons(startPathEntry, filterEntry, myWindow)

	// Initially hide progress bar container and show action button container
	toggleVisibility(true)

	// Create a content container with all the widgets
	content := container.NewVBox(
		baseDirLabel,
		container.NewVBox(startDirButton, startPathEntry),
		filterLabel,
		filterEntry,
		actionButtonContainer,
		progressBarContainer,
	)

	return content
}

func createLogScroll() fyne.CanvasObject {
	// Create a new label to show log messages
	helper.LogLabel = widget.NewLabel("")

	// Create a scrollable container for the log messages
	logScroll := container.NewScroll(helper.LogLabel)

	// Set the minimum size of the scrollable container to zero
	logScroll.SetMinSize(fyne.NewSize(0, 0))

	return logScroll
}

// listCleaning() function takes in a starting path and a filter string
// and deletes all files and folders that match the filter.
func listCleaning(startPath, filterStr string, myWindow fyne.Window) {
	// Ensure that startPath and filterStr are not empty
	if startPath == "" || filterStr == "" {
		dialog.ShowError(fmt.Errorf("Both start path and filter should be provided"), myWindow)
		return
	}

	// Clear the previous results
	helper.LogLabel.SetText("")
	toggleVisibility(false)

	// Set the progress bar value to 0
	progressBar.SetValue(0)

	// Create a regular expression from the filter string
	filter := regexp.MustCompile(filterStr + "$")

	// Get the total number of files/folders that match the filter
	total, _ := helper.Count(startPath, filter)

	// Initialize a counter for the number of files/folders deleted
	i := 0

	// Delete each file/folder that matches the filter and update the log label and progress bar
	helper.Transverse(startPath, filter, func(filename string) {
		helper.Log(fmt.Sprintf("Found %s...\n", filename))
		i++
		progressBar.SetValue(float64(i) / float64(total))
	})

	// Display a message indicating that the cleaning is complete and the number of files/folders deleted
	if i > 0 {
		dialog.ShowInformation("Cleaning complete", "All matching files/folders have been deleted", myWindow)
	} else {
		dialog.ShowInformation("Cleaning complete", "Not found any folder to be deleted", myWindow)
	}

	// Toggle visibility of the progress bar and action buttons
	toggleVisibility(true)
}

// startCleaning() takes a start path and a filter string as input
// and starts the cleaning process by deleting all matching files/folders in the start path that match the given filter.
func startCleaning(startPath, filterStr string, myWindow fyne.Window) {
	// Check if startPath and filterStr are provided
	if startPath == "" || filterStr == "" {
		dialog.ShowError(fmt.Errorf("Both start path and filter should be provided"), myWindow)
		return
	}

	// Clear the previous results and show progress bar
	helper.LogLabel.SetText("")
	progressBarContainer.Show()
	progressBar.SetValue(0)

	// Create a regex filter based on filterStr
	filter := regexp.MustCompile(filterStr + "$")

	// Count the total number of files/folders that match the filter
	total, _ := helper.Count(startPath, filter)
	i := 0

	// Delete each file/folder that matches the filter
	helper.Transverse(startPath, filter, func(filename string) {
		helper.Log(fmt.Sprintf("Deleting %s...\n", filename))
		var wg sync.WaitGroup
		wg.Add(1)
		go helper.Delete(filename, &wg)
		wg.Wait()
		i++

		// Update the progress bar
		progressBar.SetValue(float64(i) / float64(total))
	})

	// Show completion message
	if i > 0 {
		dialog.ShowInformation("Cleaning complete", "All matching files/folders have been deleted", myWindow)
	} else {
		dialog.ShowInformation("Cleaning complete", "No matching files/folders found to be deleted", myWindow)
	}

	// Hide progress bar
	progressBarContainer.Hide()
}
