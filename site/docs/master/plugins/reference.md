# Reference

View component | Description | Example
---------------|-------------|---------
[component.Text](#component-text) | Displays a string | ![newtext](newtext.png)
[component.Timestamp](#component-timestamp) | Displays Unix time in a human-readable form. Hover to see full datetime. | ![timestamp](timestamp.png)
[component.Table](#component-table) | Create a table with rows and columns | ![table](table.png)
[component.Summary](#component-summary) | Create a summary containing a header and content | ![summary](summary.png)
[component.Link](#component-link) | Create a link to an object | ![link](link.png)
[component.Quadrant](#component-quadrant) | Create a quadrant consisting of fours strings and a title. | ![quadrant](quadrant.png)
[component.Labels](#component-labels) | Create a label from a key-value pair of strings. | ![labels](labels.png)

## JSON Response Reference

For plugins written in languages besides Go, view components will need to be constructed before marshaling. Below is a reference of common components and equivalent format in JSON.

### component.Text

```go
component.NewText("component text")
```

```json
{  
   "metadata":{  
      "type":"text"
   },
   "config":{  
      "value":"component text"
   }
}
```

### component.Timestamp

```go
component.NewTimestamp(time.Unix(1559734098, 0))

```

```json
{  
   "metadata":{  
      "type":"timestamp"
   },
   "config":{  
      "timestamp":1559734098
   }
}
```

### component.Table

```go
tableCols := component.NewTableCols("Name", "Ready", "Phase", "Restarts", "Node", "Age")
component.NewTableWithRows("Pods", tableCols, []component.TableRow{
	{
		"Name":     component.NewLink("", "pod name", "/pod"),
		"Age":      component.NewTimestamp(time.Unix(1559734098, 0)),
		"Ready":    component.NewText("0/0"),
		"Restarts": component.NewText("0"),
		"Phase":    component.NewText(""),
		"Node":     component.NewText(""),
	},
})
```

```json
{  
   "metadata":{  
      "type":"table",
      "title":[  
         {  
            "metadata":{  
               "type":"text"
            },
            "config":{  
               "value":"Pods"
            }
         }
      ]
   },
   "config":{  
      "columns":[  
         {  
            "name":"Name",
            "accessor":"Name"
         },
         {  
            "name":"Ready",
            "accessor":"Ready"
         },
         {  
            "name":"Phase",
            "accessor":"Phase"
         },
         {  
            "name":"Restarts",
            "accessor":"Restarts"
         },
         {  
            "name":"Node",
            "accessor":"Node"
         },
         {  
            "name":"Age",
            "accessor":"Age"
         }
      ],
      "rows":[  
         {  
            "Age":{  
               "metadata":{  
                  "type":"timestamp"
               },
               "config":{  
                  "timestamp":1559734098
               }
            },
            "Name":{  
               "metadata":{  
                  "type":"link",
                  "title":[  
                     {  
                        "metadata":{  
                           "type":"text"
                        },
                        "config":{  
                           "value":""
                        }
                     }
                  ]
               },
               "config":{  
                  "value":"pod name",
                  "ref":"/pod"
               }
            },
            "Node":{  
               "metadata":{  
                  "type":"text"
               },
               "config":{  
                  "value":""
               }
            },
            "Phase":{  
               "metadata":{  
                  "type":"text"
               },
               "config":{  
                  "value":""
               }
            },
            "Ready":{  
               "metadata":{  
                  "type":"text"
               },
               "config":{  
                  "value":"0/0"
               }
            },
            "Restarts":{  
               "metadata":{  
                  "type":"text"
               },
               "config":{  
                  "value":"0"
               }
            }
         }
      ],
      "emptyContent":""
   }
}
```

### component.Summary

```go
sections := component.SummarySections{
	{Header: "header", Content: component.NewText("text component")},
}
component.NewSummary("summary title", sections...)
```

```json
{  
   "metadata":{  
      "type":"summary",
      "title":[  
         {  
            "metadata":{  
               "type":"text"
            },
            "config":{  
               "value":"summary title"
            }
         }
      ]
   },
   "config":{  
      "sections":[  
         {  
            "header":"header",
            "content":{  
               "metadata":{  
                  "type":"text"
               },
               "config":{  
                  "value":"text component"
               }
            }
         }
      ]
   }
}
```

### component.Link

```go
component.NewLink("link title", "pod name", "/pod")
```

```json
{  
   "metadata":{  
      "type":"link",
      "title":[  
         {  
            "metadata":{  
               "type":"text"
            },
            "config":{  
               "value":"link title"
            }
         }
      ]
   },
   "config":{  
      "value":"pod name",
      "ref":"/pod"
   }
}
```

### component.Quadrant

```go
quadrant := component.NewQuadrant("quadrant title")
quadrant.Set(component.QuadNW, "NW", "0")
quadrant.Set(component.QuadNE, "NE", "1")
quadrant.Set(component.QuadSE, "SE", "2")
quadrant.Set(component.QuadSW, "SW", "3")
```

```json
{  
   "metadata":{  
      "type":"quadrant",
      "title":[  
         {  
            "metadata":{  
               "type":"text"
            },
            "config":{  
               "value":"quadrant title"
            }
         }
      ]
   },
   "config":{  
      "nw":{  
         "value":"0",
         "label":"NW"
      },
      "ne":{  
         "value":"1",
         "label":"NE"
      },
      "se":{  
         "value":"2",
         "label":"SE"
      },
      "sw":{  
         "value":"3",
         "label":"SW"
      }
   }
}
```

### component.Labels

```go
component.NewLabels(map[string]string{"label key": "label value"})
```

```json
{  
   "metadata":{  
      "type":"labels"
   },
   "config":{  
      "labels":{  
         "label key":"label value"
      }
   }
}
```
