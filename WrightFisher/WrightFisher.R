#Check package

# Check if ggplot2 is installed, and install it if not
if (!require(ggplot2, quietly = TRUE)) {
  install.packages("ggplot2")
}

# Check if ggpubr is installed, and install it if not
if (!require(ggpubr, quietly = TRUE)) {
  install.packages("ggpubr")
}

# Check if dplyr is installed, and install it if not
if (!require(dplyr, quietly = TRUE)) {
  install.packages("dplyr")
}

# Check if viridis is installed, and install it if not
if (!require(viridis, quietly = TRUE)) {
  install.packages("viridis")
}

# Check if plotly is installed, and install it if not
if(!require(plotly,quietly = TRUE)){
  install.packages("plotly")
}

# Check if tidyverse is installed, and install it if not
if(!require(tidyverse,quietly = TRUE)){
  install.packages("tidyverse")
}

# Check if htmlwidgets is installed, and install it if not
if(!require(htmlwidgets,quietly = TRUE)){
  install.packages("htmlwidgets")
}

# Check if shiny is installed, and install it if not
if(!require(shiny,quietly = TRUE)){
  install.packages("shiny")
}

# Load required libraries
library(ggplot2)
library(ggpubr)
library(dplyr)
library(viridis)
library(plotly)
library(tidyverse)
library(htmlwidgets)
library(shiny)


# Run simulation with go
# Get the current working directory
current_dir <- getwd()

# Construct the absolute path to the Go program
go_program_path <- file.path(current_dir, "WrightFisher")

# Run the Go program from R
system(paste("go run", shQuote(go_program_path)), intern = TRUE)

# https://plotly.com/r/cumulative-animations/
accumulate_by <- function(dat, var) {
  var <- lazyeval::f_eval(var, dat)
  lvls <- plotly:::getLevels(var)
  dats <- lapply(seq_along(lvls), function(x) {
    cbind(dat[var %in% lvls[seq(1, x)], ], frame = lvls[[x]])
  })
  dplyr::bind_rows(dats)
}


# import data

# Get path
path <- getwd()
setwd(path)

# Read csv file
parameter <- read.csv('SimulationParameters.csv')
data <- read.csv('all_simulation_data.csv')

# Get the parameters
population_size <- parameter$PopulationSize
selection_coefficient <- parameter$SelectionCoefficient
start_fequency <- parameter$StartAlleleFrequency
generation_numbers <-parameter$NumGenerations
run_times <- parameter$NumberRuns


# Create a subset with average values for all generations
average_data <- data %>%
  group_by(Generations) %>%
  summarise(
    PopulationSize = mean(PopulationSize),
    SelectionCoefficient = mean(SelectionCoefficient),
    StartAlleleFrequency = mean(StartAlleleFrequency),
    Average_Num_Alleles = mean(NumAlleles),
    Average_Allele_Frequency = mean(AlleleFrequency)
  )



# make the Average Allele A Frequency Over Time Plot

# Create a line plot
lineplotA <- ggplot(average_data, aes(x = Generations, y = Average_Allele_Frequency)) +
  geom_line() +
  labs(x = "Generations", y = "Frequency of Allele A") +
  ggtitle("Average Allele A Frequency Plot")+
  theme(plot.title = element_text(hjust = 0.5))
#ylim(0, 1) #adjust the y-axis for display in plot

# display the plot (if needed)
#print(lineplotA)

# Save the plot as a png file (if needed)
#ggsave(file = "Average Allele A Frequency lineplot.png", plot = lineplotA)



# make the Average Allele a Frequency Over Time Plot

# Create a line plot
lineplota <- ggplot(average_data, aes(x = Generations, y = 1-Average_Allele_Frequency)) +
  geom_line() +
  labs(x = "Generations", y = "Frequency of Allele a") +
  ggtitle("Average Allele a Frequency Plot")+
  theme(plot.title = element_text(hjust = 0.5))
#ylim(0, 1)

# display the plot (if needed)
#print(lineplota)

# Save the plot as a png file
#ggsave(file = "Average Allele a Frequency lineplot.png", plot = lineplota)



# make average_allele_freq_plot

# make the Average allele plot that combined allele A and a
average_allele_freq_plot <- ggarrange(lineplotA, lineplota,
                                      ncol = 1, nrow = 2)

# Display the plot
print(average_allele_freq_plot)

# Save the plot as a png file
ggsave(file = "Average Allele Frequency lineplot.png", plot = average_allele_freq_plot)



# Set the data for gentype counts

# Calculate the frequency of AA, Aa, and aa in the last generation
last_generation <- max(average_data$Generations)
last_generation_data <- subset(average_data, Generations == last_generation)
last_generation_data$AA_Frequency <- (last_generation_data$Average_Allele_Frequency)**2 
last_generation_data$Aa_Frequency <- 2*last_generation_data$Average_Allele_Frequency * (1 - last_generation_data$Average_Allele_Frequency) 
last_generation_data$aa_Frequency <- (1 - last_generation_data$Average_Allele_Frequency)**2 

# Create a data frame for the bar plot
genotype_data <- data.frame(
  Genotype = c("AA", "Aa", "aa"),
  Frequency = c(last_generation_data$AA_Frequency,last_generation_data$Aa_Frequency,last_generation_data$aa_Frequency)
)



# Plot Genotype frequenct
# Create a bar plot
genotype_frequency_plot <- ggplot(genotype_data, aes(x = Genotype, y = Frequency)) +
  geom_bar(stat = "identity", fill = "skyblue") +
  ggtitle("Individual Genotype Frequency Plot of AA, Aa, and aa")+
  labs(x = "Genotype of Individuals",
       y = "Genotype Frequency") +
  theme_minimal()

# display the plot
genotype_frequency_plot

# Save the plot as a png file
ggsave(file = "Genotype Frequency Plot.png", plot = genotype_frequency_plot)



# Make Total Allele Copies Histogram

# Set the plot title
title_total <- paste("Distribution of Total Allele Copies Across Genration S(n) for s=", toString(selection_coefficient), sep = "")

# Create a ggplot object
p <- ggplot(data, aes(x = NumAlleles)) +
  geom_histogram(binwidth = 1, fill = rgb(0.2, 0.5, 0.2, 0.7), color = rgb(0.1, 0.3, 0.1, 0.7)) +
  labs(title = title_total,
       x = "Allele Count (n)",
       y = "Frequency")

# Save the plot to a png file
png_file_path <- "TotalAlleleCopiesHistogram_ggplot.png"
ggsave(png_file_path, p)

# Display the plot
print(p)



# Make dataframe for fixed and lossed plot

# Create a data frame for fixation events
fixation_data <- subset(data, AlleleFrequency == 1)

# Create a data frame for loss events
loss_data <- subset(data, AlleleFrequency == 0)

# Count the number of fixation events per generation
fixation_count <- table(fixation_data$Generations)

# Count the number of loss events per generation
loss_count <- table(loss_data$Generations)

# Create a data frame for fixation count
fixation_count_df <- data.frame(Generations = as.numeric(names(fixation_count)), Count = as.numeric(fixation_count))

# Create a data frame for loss count
loss_count_df <- data.frame(Generations = as.numeric(names(loss_count)), Count = as.numeric(loss_count))



# Create a bar plot for fixation events
fixation_bar_plot <- ggplot(fixation_count_df, aes(x = Generations, y = Count)) +
  geom_bar(stat = "identity", fill = "skyblue") +
  geom_smooth(method = "loess", se = FALSE, color = "black", linetype = "dotted") +
  labs(title = "Fixation Events Over Time",
       x = "Generations",
       y = "Count of Fixation Events")+
  ylim(0, max(loss_count_df$Count) * 1.5)

# Save the plot as a png file
#ggsave(file = "Fixation_Count_Bar_Plot.png", plot = fixation_bar_plot)

# Display the plot
#print(fixation_bar_plot)



# Create a bar plot for loss events
loss_bar_plot <- ggplot(loss_count_df, aes(x = Generations, y = Count)) +
  geom_bar(stat = "identity", fill = "coral") +
  geom_smooth(method = "loess", se = FALSE, color = "black", linetype = "dotted") +
  labs(title = "Loss Events Over Time",
       x = "Generations",
       y = "Count of Loss Events")+
  ylim(0, max(loss_count_df$Count) * 1.5)

# Save the plot as a png file
#ggsave(file = "Loss_Count_Bar_Plot.png", plot = loss_bar_plot)

# Display the plot
#print(loss_bar_plot)


# make Fixation and loss bar plots that conbain fix and loss plot
fix_loss_bar_polt <- ggarrange(fixation_bar_plot, loss_bar_plot,
                               ncol = 2, nrow = 1)

print(fix_loss_bar_polt)
# Save the plot as a png file
ggsave(file = "Fixation and loss bar plots.png", plot = fix_loss_bar_polt)



# Count the maximum allele number reached for each simulation
max_allele_counts <- tapply(data$NumAlleles, data$Generations, max)

# Create a data frame for maximum allele counts
max_allele_df <- data.frame(Generations = as.numeric(names(max_allele_counts)), MaxAllele = as.numeric(max_allele_counts))

# Create a histogram for the distribution of maximum allele number reached
max_allele_plot <- ggplot(max_allele_df, aes(x = MaxAllele)) +
  geom_histogram(binwidth = 1, fill = "darkgreen", color = "black", alpha = 0.7) +
  labs(title = paste("Distribution of Maximum Allele Number Reached (Mi) for s=",toString(selection_coefficient), sep = ""),
       x = "Maximum Allele Number (Mi)",
       y = "Number of Simulations") +
  theme_minimal()

# Save the plot as a png file
ggsave(file = "Max Allele Histogram.png", plot = max_allele_plot)

# Display the plot
print(max_allele_plot)



# Create a combined plot with density on top and heatmap at the bottom
combined_plot <- ggplot(data, aes(x = Generations, y = NumAlleles)) +
  geom_tile(aes(fill = NumAlleles), alpha = 0.7) +  # Heatmap at the bottom
  geom_density_2d(aes(fill = after_stat(level)), contour = FALSE) +  # Density plot on top
  scale_fill_viridis(name = "NumAlleles", guide = "legend") +  # Adjust the color scale as needed
  labs(title = paste("Combined Plot for", run_times, "Simulations (s =", selection_coefficient, ")"),
       x = "Time (Generations)",
       y = "Allele Count",
       fill = "") +  # Remove legend title
  theme_minimal()

# Save the combined plot as a png file
ggsave(file = "Combined_Plot.png", plot = combined_plot)

# Display the combined plot
print(combined_plot)



# make a merged dataframe that conatin the generation allele numbe, allele frquency, fixed count and loss count

merged_df <- subset(average_data, select = c(Generations, Average_Num_Alleles, Average_Allele_Frequency))

merged_df <- merged_df %>%
  full_join(fixation_count_df, by = 'Generations') %>%
  full_join(loss_count_df, by = 'Generations') %>%
  mutate(FixedNumber = Count.x, LossedNumber = Count.y) %>%
  select(-Count.x, -Count.y)

merged_df$Fix_loss_ratio = merged_df$FixedNumber/merged_df$LossedNumber



# Alelle fixed and lossed ratio change Plot
lineplot_flratio <- ggplot(merged_df, aes(x = Generations, y = Fix_loss_ratio)) +
  geom_line() +
  labs(x = "Generations", y = "The Ration of Alelle fixed and lossed") +
  ggtitle("Alelle fixed and lossed ratio Plot")+
  theme(plot.title = element_text(hjust = 0.5))

lineplot_flratio

# Save the plot as a png file
ggsave(file = "Alelle fixed and lossed ratio Plot.png", plot = lineplot_flratio)



# Make the cumulative line plot animation

# Calculate the frequency of AA, Aa, and aa in all generations
df_genotype_num <- subset(average_data, select = c(Generations,PopulationSize, Average_Allele_Frequency)) 
df_genotype_num$AA <- (df_genotype_num$Average_Allele_Frequency)**2 * df_genotype_num$PopulationSize
df_genotype_num$Aa <- 2*df_genotype_num$Average_Allele_Frequency * (1 - df_genotype_num$Average_Allele_Frequency) * df_genotype_num$PopulationSize
df_genotype_num$aa <- (1 - df_genotype_num$Average_Allele_Frequency)**2 * df_genotype_num$PopulationSize

df_animation <- df_genotype_num %>%
  # Reshape the data to long format
  pivot_longer(cols = c(AA, Aa, aa), names_to = "Genotype", values_to = "Allele_Numbers") %>%
  # Select only the relevant columns
  select(Generations, Allele_Numbers, Genotype)

fig_genotype <- df_animation %>% accumulate_by(~Generations)


fig_genotype <- fig_genotype %>%
  plot_ly(
    x = ~Generations, 
    y = ~Allele_Numbers,
    split = ~Genotype,
    frame = ~frame, 
    type = 'scatter',
    mode = 'lines', 
    line = list(simplyfy = F)
  )
fig_genotype <- fig_genotype %>% layout(
  xaxis = list(
    title = "Generations",
    zeroline = F
  ),
  yaxis = list(
    title = "Average number of each genotype",
    zeroline = F
  )
) 
fig_genotype <- fig_genotype %>% animation_opts(
  frame = 100, 
  transition = 0, 
  redraw = FALSE
)
fig_genotype <- fig_genotype %>% animation_slider(
  hide = T
)
fig_genotype <- fig_genotype %>% animation_button(
  x = 1, xanchor = "right", y = 0, yanchor = "bottom"
)

fig_genotype

# Save the plot as an HTML file
htmlwidgets::saveWidget(fig_genotype, file = "fig_genotype.html")


# Allele A frequency in 5 runs

# Create a new data frame with the first five runs
data_five_runs <- data %>% slice(1:(5*generation_numbers))

data_five_runs <- subset(data_five_runs,select = c(Generations,AlleleFrequency))

# Add a new column "Runs" based on row index
data_five_runs <- data_five_runs %>%
  mutate(Runs = rep(1:(n() %/% generation_numbers), each = generation_numbers, length.out = n()))


fig_5runs <- data_five_runs %>% accumulate_by(~Generations)

fig_5runs <- fig_5runs %>%
  plot_ly(
    x = ~Generations, 
    y = ~AlleleFrequency,
    split = ~Runs,
    frame = ~frame, 
    type = 'scatter',
    mode = 'lines', 
    line = list(simplyfy = F)
  )
fig_5runs <- fig_5runs %>% layout(
  xaxis = list(
    title = "Generations",
    zeroline = F
  ),
  yaxis = list(
    title = "Allele Frequency of A",
    zeroline = F
  ),
  legend = list(
    title = "Runs"  # Set the legend title for the entire legend box
  )
) 
fig_5runs <- fig_5runs %>% animation_opts(
  frame = 100, 
  transition = 0, 
  redraw = FALSE
)
fig_5runs <- fig_5runs %>% animation_slider(
  hide = T
)
fig_5runs <- fig_5runs %>% animation_button(
  x = 1, xanchor = "right", y = 0, yanchor = "bottom"
)

fig_5runs

# Save the plot as an HTML file
htmlwidgets::saveWidget(fig_5runs, file = "fig_5runs.html")



