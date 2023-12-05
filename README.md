# DXDM_project
 02601_final_project

Wright-Fisher model simulation is a program that uses Go language to simulate and R language for visualization.

Before starting the simulation, PLEASE make sure the GO workplace is all setup, also some required packages must be added to the Go path: gonum.org/v1/gonum/stat/distuv
To install this package, you can use this code:
go install gonum.org/v1/gonum/stat/distuv@latest



The other packages in both Go and R can be automatically installed if they do not exist.
If the package load fails, please try to manually add packages. Here are all the packages required:
Go: 	"encoding/csv", "fmt",		"log", "time", "os",	"path/filepath", "gonum.org/v1/gonum/stat/distuv"
R: "ggplot2", "ggpubr", "dplyr", "viridis", "plotly", "tidyverse", "htmlwidgets","shiny", and "shinyjs" 

Open Rstudio and file "server.R" library all packages first, and type the code "runApp()" in Console or click the "Run App" on the top-right corner can start the app.

**** PLEASE!!! library all packages first before running the R shiny app, if you see "An error has occurred! could not find function "plotlyOutput"" this error code, you are not library packages. ****

**** The WrightFisher.R is the raw cod, you can get a specific plot with specific parameters using that file. You also can use the Go to simulation and output to csv file, the code to start is "./WrightFisherSimulation populationSize selectCoefficent startFrequency generationNumber runTimes". For example: ./WrightFisherSimulation 200 0 0.5 100 100" ****
