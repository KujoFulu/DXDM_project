library(shiny)
library(shinyMatrix)
library(ggplot2)
library(reshape2)
library(dplyr)
library(tidyr)
library(akima)

# Define server logic
server <- function(input, output, session) {
  # Run Simulation and plotting
  runSimulation <- function() {
    # Trigger a notice
    showNotification("Starting Lotka-Volterra simulation...", type = "message")
    
    numSpecies <- input$numSpecies
    
    # Get initial populations for each species
    initPop <-
      sapply(1:numSpecies, function(i)
        input[[paste("initPop", i, sep = "")]])
    
    # Get interaction matrix and collapse it into a 1D vector
    interactionMatrix <- c(input$interactionMatrix)
    
    # Get growth/death rates for each species
    rate <-
      sapply(1:numSpecies, function(i)
        input[[paste("rate", i, sep = "")]])
    
    # Set the working directory to the location of the Go program
    setwd("/Users/graceya/go/src/LotkaVolterra")
    current_dir <- getwd()
    go_program_path <-
      file.path(current_dir, "LVSimulation")
    
    # Construct the command with arguments for the Go program
    command <- paste(
      "go run",
      shQuote(go_program_path),
      numSpecies,
      paste(initPop, collapse = " "),
      paste(interactionMatrix, collapse = " "),
      paste(rate, collapse = " "),
      sep = " "
    )
    
    # Run the simulation
    simulationOutput <- system(command, intern = TRUE, wait = TRUE)
    print(simulationOutput)  # This will print the output of the command, including errors
    
    # Trigger a notice
    showNotification("Simulation completed!", type = "message")
    
    # Load the simulation data
    # Update the file path if needed
    simulationData <- read.csv("LVSimulation/output/test.csv")
    
    # Display the simulation data as a table
    output$simulationData <- renderTable({
      head(simulationData, 50)
    })
    
    # Display the gif
    output$animated_gif <- renderImage({
      # Return a list containing the filename
      list(src = "LVSimulation/output/test.out.gif",
           contentType = "image/gif",
           alt = "This is alternate text")
    }, deleteFile = FALSE)
    
    # Display population over time plot using a loop depending on the number of species
    output$time_plot <- renderPlot({
      simulationData <- read.csv("LVSimulation/output/test.csv")
      # Transform the data to a long format
      long_data <-
        melt(simulationData,
             id.vars = "Generation",
             variable.name = "Species")
      long_data$Species <-
        factor(long_data$Species)
      # Plot
      ggplot(long_data,
             aes(
               x = Generation,
               y = value,
               group = Species,
               color = Species
             )) +
        geom_line(size = 1.5) +
        scale_color_manual(
          values = c(
            "#E41A1C",
            "#377EB8",
            "#4DAF4A",
            "#984EA3",
            "#FF7F00",
            "#FFFF33",
            "#A65628",
            "#F781BF",
            "#999999"
          )
        ) +
        labs(x = "Number of Generations",
             y = "Population Size (in thousands)",
             color = "Species") +
        theme_minimal() +
        theme(
          legend.position = "bottom",
          # Move legend to bottom
          legend.text = element_text(size = 16),
          legend.title = element_text(size = 16),
          axis.text.x = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.text.y = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.title.x = element_text(size = 18, color = "black"),
          axis.title.y = element_text(size = 18, color = "black")
        )
      
    })
    
    output$phase_plane_plot <- renderPlot({
      # Create the phase plane plot
      ggplot(simulationData, aes(x = `Species.0`, y = `Species.1`)) +
        geom_point(alpha = 0.5, size =0.3) +  # Use points with some transparency
        geom_path(color = '#6A0DAD', size = 0.5) +  # Connect points with a path
        labs(x = "Population of Species 0 (in thousands)",
             y = "Population of Species 1 (in thousands)") +
        theme_minimal() +
        theme(
          axis.text.x = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.text.y = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.title.x = element_text(size = 18, color = "black"),
          axis.title.y = element_text(size = 18, color = "black")
        )
    })
    
    # Calculate the sequence of rows to select
    rows_to_select <- seq(1, nrow(simulationData), by = 2499)
    # Subset simulationData using the calculated sequence
    simulationDataSubset <- simulationData[rows_to_select, ]
    
    # Plot the vector field
    output$vector_field_plot <- renderPlot({
      # Calculate differences (gradients) for Species 0 and Species 1
      simulationDataSubset <- simulationDataSubset %>%
        mutate(dSpecies0 = c(NA, diff(`Species.0`)),
               dSpecies1 = c(NA, diff(`Species.1`)))
      
      # Define the grid for the vector field
      nb_points <- 20
      x_seq <-
        seq(0, max(simulationDataSubset$`Species.0`), length.out = nb_points)
      y_seq <-
        seq(0, max(simulationDataSubset$`Species.1`), length.out = nb_points)
      
      # Create a meshgrid of points
      grid <- expand.grid(x = x_seq, y = y_seq)
      
      # Define your system of equations here
      # Define the Lotka-Volterra system of equations
      # Ensure the specific elements you want to use are numeric
      alpha <- as.numeric(rate[1])
      beta <- as.numeric(interactionMatrix[numSpecies + 1])
      gamma <- as.numeric(rate[2])
      delta <- as.numeric(interactionMatrix[2])
      
      if (!is.numeric(alpha) || !is.numeric(beta) ||
          !is.numeric(gamma) || !is.numeric(delta)) {
        stop("One of the required elements is not numeric")
      }
      
      dX_dt <- function(X, Y) {
        list(DX = alpha * X + beta * X * Y,
             DY = delta * X * Y + gamma * Y)
      }
      
      # Calculate the direction at each point in the grid
      vectorField <- mapply(dX_dt, grid$x, grid$y, SIMPLIFY = FALSE)
      vectorField <- do.call(rbind, vectorField)
      
      # Combine the grid and vectorField into a data frame for ggplot
      vectorFieldData <- cbind(grid, vectorField)
      
      # Ensure DX and DY are numeric
      vectorFieldData$DX <- as.numeric(vectorFieldData$DX)
      vectorFieldData$DY <- as.numeric(vectorFieldData$DY)
      
      # Check for NA values in DX or DY
      if (any(is.na(vectorFieldData$DX)) ||
          any(is.na(vectorFieldData$DY))) {
        stop("NA values found in DX or DY")
      }
      
      # Compute xend and yend
      vectorFieldData$xend <-
        vectorFieldData$x + 0.1 * vectorFieldData$DX
      vectorFieldData$yend <-
        vectorFieldData$y + 0.1 * vectorFieldData$DY
      
      ggplot() +
        geom_segment(
          data = vectorFieldData,
          aes(
            x = x,
            y = y,
            xend = xend,
            yend = yend,
            color = "#D7B9D5"
          ),
          arrow = arrow(type = "closed", length = unit(0.1, "inches")),
          size = 0.5
        ) +
        labs(x = "Population of Species 0 (in thousands)", y = "Population of Species 1 (in thousands)") +
        theme_minimal() +
        theme(
          axis.text.x = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.text.y = element_text(
            size = 16,
            color = "black",
            angle = 0,
            hjust = 0.5,
            vjust = 0.5
          ),
          axis.title.x = element_text(size = 18, color = "black"),
          axis.title.y = element_text(size = 18, color = "black"),
          legend.position = "none"
        )
    })
    
    # Plot the contour plot
    output$contour_plot <- renderPlot({
      # Ensure Species.0 is numeric
      simulationDataSubset$`Species.0` <-
        as.numeric(as.character(simulationDataSubset$`Species.0`))
      
      # Interpolate data onto a regular grid, handling duplicate points
      interpolated <-
        with(
          simulationDataSubset,
          interp(
            x = `Species.1`,
            y = `Species.2`,
            z = `Species.0`,
            duplicate = "mean"
          )
        )
      
      # Convert to data frame for ggplot
      contourData <- as.data.frame(interp2xyz(interpolated))
      
      # Create the contour plot
      ggplot(contourData, aes(x = x, y = y, z = z)) +
        geom_contour_filled() +
        labs(x = "Population of Species 1 (in thousands)", 
             y = "Population of Species 2 (in thousands)",
             fill = "Population of Species 0") +
        theme_minimal() +
        theme(
          legend.position = "bottom",
          # Move legend to bottom
          legend.text = element_text(size = 16),
          legend.title = element_text(size = 16),
          axis.text.x = element_text(size = 16),
          axis.text.y = element_text(size = 16),
          axis.title.x = element_text(size = 18),
          axis.title.y = element_text(size = 18)
        ) +
        ggtitle("Contour Plot: Species 0 Levels on Species 1 and 2 Plane")
    })
    
    # Trigger a notice
    showNotification("Plots Generated!", type = "default")
  }
  
  # Observing the button click
  observeEvent(input$original, {
    # Define defaults
    originalNumSpecies <- 3
    originalInitPop <- c(50, 10, 5)
    originalInteractionMatrix <- matrix(
      c(0,-0.04,-0.04,
        0.04, 0,-0.02,
        0.02, 0.04, 0),
      nrow = originalNumSpecies,
      ncol = originalNumSpecies,
      byrow = TRUE
    )
    originalRate <- c(0.25,-0.5,-0.5)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = originalNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:originalNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = originalInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = originalRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = originalInteractionMatrix,
        rows = list(names = paste("Species", 1:originalNumSpecies)),
        cols = list(names = paste("Species", 1:originalNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
    
  })
  
  # observe "chaotic dynamic" button
  observeEvent(input$chaotic, {
    # Define defaults
    chaoticNumSpecies <- 4
    chaoticInitPop <- c(0.1, 0.8, 0.3, 0.5)
    chaoticInteractionMatrix <- matrix(
      c(
        -1,
        -1.09,
        -1.52,
        0,
        0,
        -0.72,
        -0.3168,
        -0.9792,-3.5649,
        0,
        -1.53,
        -0.7191,-1.5367,
        -0.6477,
        -0.4445,
        -1.27
      ),
      nrow = chaoticNumSpecies,
      ncol = chaoticNumSpecies,
      byrow = TRUE
    )
    chaoticRate <- c(1, 0.72, 1.53, 1.27)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = chaoticNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:chaoticNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = chaoticInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = chaoticRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = chaoticInteractionMatrix,
        rows = list(names = paste("Species", 1:chaoticNumSpecies)),
        cols = list(names = paste("Species", 1:chaoticNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
    
  })
  
  # observe "extinction of One Species" button
  observeEvent(input$extinctionOne, {
    # Define defaults
    extinctionNumSpecies <- 3
    extinctionInitPop <- c(0.1, 0.8, 0.3)
    extinctionInteractionMatrix <- matrix(
      c(-2,-1,-1,-1,-1,-2,-2.6,-1.6,-3),
      nrow = extinctionNumSpecies,
      ncol = extinctionNumSpecies,
      byrow = TRUE
    )
    extinctionRate <- c(3, 4, 7.2)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = extinctionNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:extinctionNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = extinctionInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = extinctionRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = extinctionInteractionMatrix,
        rows = list(names = paste("Species", 1:extinctionNumSpecies)),
        cols = list(names = paste("Species", 1:extinctionNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
    
  })
  
  # observe "extinction of Two Species" button
  observeEvent(input$extinctionTwo, {
    # Define defaults
    extinction2NumSpecies <- 3
    extinction2InitPop <- c(0.1, 0.8, 0.3)
    extinction2InteractionMatrix <- matrix(
      c(-0.1,-1,-0.1,-1,-0.1,-2,-2.6,-0.6,-3),
      nrow = extinction2NumSpecies,
      ncol = extinction2NumSpecies,
      byrow = TRUE
    )
    extinction2Rate <- c(3, 4, 7.2)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = extinction2NumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:extinction2NumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = extinction2InitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = extinction2Rate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = extinction2InteractionMatrix,
        rows = list(names = paste("Species", 1:extinction2NumSpecies)),
        cols = list(names = paste("Species", 1:extinction2NumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
    
  })
  
  # observe "Limit Cycle" button
  observeEvent(input$limit, {
    # Define defaults
    limitNumSpecies <- 3
    limitInitPop <- c(0.1, 0.8, 0.3)
    limitInteractionMatrix <- matrix(
      c(-0.5,-1, 0, 0,-1,-2,-2.6,-1.6,-3),
      nrow = limitNumSpecies,
      ncol = limitNumSpecies,
      byrow = TRUE
    )
    limitRate <- c(3, 4, 7.2)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = limitNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:limitNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = limitInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = limitRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = limitInteractionMatrix,
        rows = list(names = paste("Species", 1:limitNumSpecies)),
        cols = list(names = paste("Species", 1:limitNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
  })
  
  # observe "lynx" button
  observeEvent(input$lynx, {
    # Define defaults
    lynxNumSpecies <- 2
    lynxInitPop <- c(3.49, 0.38)
    lynxInteractionMatrix <- matrix(
      c(0,-0.02482, 0.02756, 0),
      nrow = lynxNumSpecies,
      ncol = lynxNumSpecies,
      byrow = TRUE
    )
    lynxRate <- c(0.48,-0.93)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = lynxNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:lynxNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = lynxInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = lynxRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = lynxInteractionMatrix,
        rows = list(names = paste("Species", 1:lynxNumSpecies)),
        cols = list(names = paste("Species", 1:lynxNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
  })
  
  # observe "stable" button
  observeEvent(input$stable, {
    # Define defaults
    stableNumSpecies <- 3
    stableInitPop <- c(0.1, 0.8, 0.3)
    stableInteractionMatrix <- matrix(
      c(-2,-1, 0, 0,-1,-2,-2.6,-1.6,-3),
      nrow = stableNumSpecies,
      ncol = stableNumSpecies,
      byrow = TRUE
    )
    stableRate <- c(3, 4, 7.2)
    
    # Update numSpecies input to default
    updateNumericInput(session, "numSpecies", value = stableNumSpecies)
    
    # Now update the dynamic UI inputs
    # Since these inputs are dynamic, you must recreate them with the correct default values
    output$dynamicUI <- renderUI({
      uiElements <- list()
      
      # Add sliders for initial population for each species with default values
      for (i in 1:stableNumSpecies) {
        uiElements[[length(uiElements) + 1]] <- sliderInput(
          inputId = paste("initPop", i, sep = ""),
          label = paste("Initial Population for Species", i, ":"),
          min = 0,
          max = 100,
          value = stableInitPop[i],
          step = 0.1
        )
        
        # Add numeric inputs for growth/death rate for each species with default values
        uiElements[[length(uiElements) + 1]] <- numericInput(
          inputId = paste("rate", i, sep = ""),
          label = paste("Growth/Death Rate for Species", i, ":"),
          min = -10,
          max = 10,
          value = stableRate[i],
          step = 0.01
        )
      }
      
      # Add interaction matrix with default values
      uiElements[[length(uiElements) + 1]] <- matrixInput(
        "interactionMatrix",
        label = "Interaction Matrix:",
        value = stableInteractionMatrix,
        rows = list(names = paste("Species", 1:stableNumSpecies)),
        cols = list(names = paste("Species", 1:stableNumSpecies))
      )
      
      do.call(tagList, uiElements)
    })
  })
  
  # Reactive expression to create sliders and a single rate input based on number of species
  output$dynamicUI <- renderUI({
    species <- input$numSpecies
    
    # Create a list to store UI elements
    uiElements <- list()
    
    # Add sliders for initial population and a single numeric input for growth/death rate for each species
    for (i in 1:species) {
      # Initial Population Slider
      uiElements[[length(uiElements) + 1]] <- sliderInput(
        inputId = paste("initPop", i, sep = ""),
        label = paste("Initial Population for Species", i, ":"),
        min = 0,
        max = 100,
        value = 50,
        step = 0.1
      )
      
      # Growth/Death Rate Numeric Input
      uiElements[[length(uiElements) + 1]] <- numericInput(
        inputId = paste("rate", i, sep = ""),
        label = paste("Growth/Death Rate for Species", i, ":"),
        min = -10,
        max = 10,
        value = 0.25,
        step = 0.01
      )
    }
    
    # Add interaction matrix
    uiElements[[length(uiElements) + 1]] <- matrixInput(
      "interactionMatrix",
      label = "Interaction Matrix:",
      value = matrix(0, species, species),
      rows = list(names = paste("Species", 1:species)),
      cols = list(names = paste("Species", 1:species))
    )
    
    do.call(tagList, uiElements)
  })

  # Define the observeEvent for the "Run Simulation" button
  observeEvent(input$runSimulation, {
    runSimulation()
  })
  
  # Render the Save All Plots button
  observeEvent(input$savePlots, {
    
    # Trigger a notice
    showNotification("Saving Plots",type = "default" )
    
    # Save the Average Allele Frequency Plot
    ggsave(file = "vector_field_plot.png", plot =vector_field_plot )
    
    # Trigger a notice
    showNotification("All plots have been saved successfully!",type = "message" )
  })
  
}
