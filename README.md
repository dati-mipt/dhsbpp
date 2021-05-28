dhsbpp
================================================================

Summary
-------

dhsbpp is a tool for solving Dynamical Hierarchical Structured Bin Packing Problem. It supports tuning of various parameters
related to algorithm.  

Installation
----------------
    go get -u github.com/dati-mipt/dhsbpp
(the same command works for updating)

Example usage
---------------

    go run github.com/dati-mipt/dhsbpp/main -max_capacity=100000 -init_days=10 -dataset=australia -separate=root -algorithm=first_fit

Datasets
---------------
[Datasets](https://github.com/ffuf/ffuf/blob/master/LICENSE) for the algorithm include:
- Tree topology (`datasets/%DATASETNAME%/ChildParent.csv`) 
- Node weights for consecutive epochs (`datasets/%DATASETNAME%/WeightsPerEpoch.csv`)


Command-line options
---------------------
```
-max_capacity=N, must be specified
    Maximum bin size

-init_epochs=N, must be specified
    Number of epochs for initial distribution of node weights.
    
-dataset=<folder>, must be specified
    Valid dataset for algorithm
   
-algorithm=first_fit/greedy, must be specified
    Specifies the packing algorithm.
    
-separate=max_child/root, must be specified
    Specifies the way of separating nodes in packing algorithm.
    
-AF=N, default: 0.6
    0 < N < 1
    Indicates the maximum allowed usage of a single bin during (re)allocation.
    If a single bin has a maximum capacity of MAX, the bin packing algorithms 
    will be invoked with AF × MAX for the size of a single bin.
    
    Note: AF + RD < 1

-RD=N, default: 0.2
    0 < N < 1
    Used for reallocations. The trigger for reallocation is the overload 
    or underload of a bin.
    
    Note: AF + RD < 1
```


Output images
---------------------

The tool generates images of consecutive tree distributions.

*TODO Upload images*