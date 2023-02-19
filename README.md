# fairytale
Develop and document Hypp components

## Package dependency graph

```mermaid
flowchart TD

subgraph internal
    dispatch
    component
    console
    driver
end


book --> dispatch
book --> component
book --> fairytale

dispatch --> driver
dispatch --> fairytale

component --> dispatch
component --> fairytale

control --> fairytale
control --> dispatch

fairytale --> console
```
