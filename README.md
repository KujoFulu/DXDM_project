# DXDM_project
 02601_final_project

# Wright-Fisher Model
Wright-Fisher model simulation is a program that uses Go language to simulate and R language for visualization.

Before starting the simulation, PLEASE make sure the GO workplace is all set, also some required packages must be added to the Go path: gonum.org/v1/gonum/stat/distuv
To install this package, you can use this code:
"go install gonum.org/v1/gonum/stat/distuv@latest"

The other packages in both Go and R can be automatically installed if they do not exist.

If the package load fails, please try to add packages manually. Here are all the packages required:
Go: 	"encoding/csv", "fmt",		"log", "time", "os",	"path/filepath", "gonum.org/v1/gonum/stat/distuv"
R: "ggplot2", "ggpubr", "dplyr", "viridis", "plotly", "tidyverse", "htmlwidgets","shiny", and "shinyjs" 

Open Rstudio and file "server.R" library all packages first, and type the code "runApp()" in the Console or click the "Run App" on the top-right corner to start the app.

 **** PLEASE!!! library all packages first before running the R shiny app. If you see "An error has occurred! could not find function "plotlyOutput" this error code, you are not library packages. ****

 **** The WrightFisher.R is the raw cod; you can get a specific plot with specific parameters using that file. You also can use the Go to simulation and output to a CSV file. The code to start is "./WrightFisherSimulation populationSize selectCoefficent startFrequency generationNumber runTimes". For example "./WrightFisherSimulation 200 0 0.5 100 100" ****

# Lotka-Volterra Model
This model is a program uses Go language to simulate, Python and R for visulization.

Before starting the simulation, PLEASE make sure the GO workplace is all set, also some required packages must be added to the Go path: gonum.org/v1/gonum/mat
To install this package, you can use this code:
"go install gonum.org/v1/gonum/mat"

The packages for Pyhton code to draw plots are "pandas", "matplotlib.pyplot" and "numpy", you can use conda to install them.

If you want to seperately run the Go simulation, you can go into  LotkaVolterra/LVSimulation folder, do the following two command lines:
go build
./LVSimulation 3 50.0 10.0 5.0 0 0.04 0.02 -0.04 0 0.04 -0.04 -0.02 0 0.25 -0.5 -0.5