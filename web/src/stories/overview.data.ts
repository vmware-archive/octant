export const big_data = {
  "metadata": {
    "type": "list",
    "title": [{"metadata": {"type": "text"}, "config": {"value": "Edit JSON data in Knobs to modify this view"}}]
  },
  "config": {
    "items": [{
      "metadata": {
        "type": "table",
        "title": [{"metadata": {"type": "text"}, "config": {"value": "Services"}}]
      },
      "config": {
        "columns": [{"name": "Name", "accessor": "Name"}, {
          "name": "Labels",
          "accessor": "Labels"
        }, {"name": "Type", "accessor": "Type"}, {
          "name": "Cluster IP",
          "accessor": "Cluster IP"
        }, {"name": "External IP", "accessor": "External IP"}, {
          "name": "Ports",
          "accessor": "Ports"
        }, {"name": "Age", "accessor": "Age"}, {"name": "Selector", "accessor": "Selector"}],
        "rows": [{
          "Age": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1588715238}},
          "Cluster IP": {"metadata": {"type": "text"}, "config": {"value": "10.96.0.1"}},
          "External IP": {"metadata": {"type": "text"}, "config": {"value": "<none>"}},
          "Labels": {
            "metadata": {"type": "labels"},
            "config": {"labels": {"component": "apiserver", "provider": "kubernetes"}}
          },
          "Name": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "kubernetes",
              "ref": "/overview/namespace/default/discovery-and-load-balancing/services/kubernetes",
              "status": 1,
              "statusDetail": {
                "metadata": {"type": "list"},
                "config": {"items": [{"metadata": {"type": "text"}, "config": {"value": "Service is OK"}}]}
              }
            }
          },
          "Ports": {"metadata": {"type": "text"}, "config": {"value": "443/TCP"}},
          "Selector": {"metadata": {"type": "selectors"}, "config": {"selectors": null}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "ClusterIP"}},
          "_action": {
            "metadata": {"type": "gridActions"},
            "config": {
              "actions": [{
                "name": "Delete",
                "actionPath": "action.octant.dev/deleteObject",
                "payload": {"apiVersion": "v1", "kind": "Service", "name": "kubernetes", "namespace": "default"},
                "confirmation": {
                  "title": "Delete Service",
                  "body": "Are you sure you want to delete *Service* **kubernetes**? This action is permanent and cannot be recovered."
                },
                "type": "danger"
              }]
            }
          }
        }],
        "emptyContent": "We couldn't find any services!",
        "loading": false,
        "filters": {}
      }
    }, {
      "metadata": {"type": "table", "title": [{"metadata": {"type": "text"}, "config": {"value": "Secrets"}}]},
      "config": {
        "columns": [{"name": "Name", "accessor": "Name"}, {
          "name": "Labels",
          "accessor": "Labels"
        }, {"name": "Type", "accessor": "Type"}, {"name": "Data", "accessor": "Data"}, {
          "name": "Age",
          "accessor": "Age"
        }],
        "rows": [{
          "Age": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1588715247}},
          "Data": {"metadata": {"type": "text"}, "config": {"value": "3"}},
          "Labels": {"metadata": {"type": "labels"}, "config": {"labels": {}}},
          "Name": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "default-token-k4mp4",
              "ref": "/overview/namespace/default/config-and-storage/secrets/default-token-k4mp4",
              "status": 1,
              "statusDetail": {
                "metadata": {"type": "list"},
                "config": {"items": [{"metadata": {"type": "text"}, "config": {"value": "v1 Secret is OK"}}]}
              }
            }
          },
          "Type": {"metadata": {"type": "text"}, "config": {"value": "kubernetes.io/service-account-token"}},
          "_action": {
            "metadata": {"type": "gridActions"},
            "config": {
              "actions": [{
                "name": "Delete",
                "actionPath": "action.octant.dev/deleteObject",
                "payload": {
                  "apiVersion": "v1",
                  "kind": "Secret",
                  "name": "default-token-k4mp4",
                  "namespace": "default"
                },
                "confirmation": {
                  "title": "Delete Secret",
                  "body": "Are you sure you want to delete *Secret* **default-token-k4mp4**? This action is permanent and cannot be recovered."
                },
                "type": "danger"
              }]
            }
          }
        }],
        "emptyContent": "We couldn't find any secrets!",
        "loading": false,
        "filters": {}
      }
    }, {
      "metadata": {
        "type": "table",
        "title": [{"metadata": {"type": "text"}, "config": {"value": "Service Accounts"}}]
      },
      "config": {
        "columns": [{"name": "Name", "accessor": "Name"}, {
          "name": "Labels",
          "accessor": "Labels"
        }, {"name": "Secrets", "accessor": "Secrets"}, {"name": "Age", "accessor": "Age"}],
        "rows": [{
          "Age": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1588715246}},
          "Labels": {"metadata": {"type": "labels"}, "config": {"labels": {}}},
          "Name": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "default",
              "ref": "/overview/namespace/default/config-and-storage/service-accounts/default",
              "status": 1,
              "statusDetail": {
                "metadata": {"type": "list"},
                "config": {
                  "items": [{
                    "metadata": {"type": "text"},
                    "config": {"value": "v1 ServiceAccount is OK"}
                  }]
                }
              }
            }
          },
          "Secrets": {"metadata": {"type": "text"}, "config": {"value": "1"}},
          "_action": {
            "metadata": {"type": "gridActions"},
            "config": {
              "actions": [{
                "name": "Delete",
                "actionPath": "action.octant.dev/deleteObject",
                "payload": {
                  "apiVersion": "v1",
                  "kind": "ServiceAccount",
                  "name": "default",
                  "namespace": "default"
                },
                "confirmation": {
                  "title": "Delete ServiceAccount",
                  "body": "Are you sure you want to delete *ServiceAccount* **default**? This action is permanent and cannot be recovered."
                },
                "type": "danger"
              }]
            }
          }
        }],
        "emptyContent": "We couldn't find any service accounts!",
        "loading": false,
        "filters": {}
      }
    }, {
      "metadata": {"type": "table", "title": [{"metadata": {"type": "text"}, "config": {"value": "Events"}}]},
      "config": {
        "columns": [{"name": "Kind", "accessor": "Kind"}, {
          "name": "Message",
          "accessor": "Message"
        }, {"name": "Reason", "accessor": "Reason"}, {"name": "Type", "accessor": "Type"}, {
          "name": "First Seen",
          "accessor": "First Seen"
        }, {"name": "Last Seen", "accessor": "Last Seen"}],
        "rows": [{
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342696}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (1)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342696}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Node minikube event: Registered Node minikube in Controller",
              "ref": "/overview/namespace/default/events/minikube.1619233ef88b9048"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "RegisteredNode"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342680}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (1)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342680}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Starting kube-proxy.",
              "ref": "/overview/namespace/default/events/minikube.1619233b5373d195"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "Starting"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342654}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (7)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342655}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Node minikube status is now: NodeHasSufficientPID",
              "ref": "/overview/namespace/default/events/minikube.16192335109a3bee"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "NodeHasSufficientPID"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342654}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (8)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342655}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Node minikube status is now: NodeHasNoDiskPressure",
              "ref": "/overview/namespace/default/events/minikube.1619233510993cda"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "NodeHasNoDiskPressure"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342654}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (8)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342655}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Node minikube status is now: NodeHasSufficientMemory",
              "ref": "/overview/namespace/default/events/minikube.161923351097f4ba"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "NodeHasSufficientMemory"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342654}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (1)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342654}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Updated Node Allocatable limit across pods",
              "ref": "/overview/namespace/default/events/minikube.161923352af5316a"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "NodeAllocatableEnforced"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342653}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (1)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592342653}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Starting kubelet.",
              "ref": "/overview/namespace/default/events/minikube.16192334f653ead2"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "Starting"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }, {
          "First Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592272167}},
          "Kind": {"metadata": {"type": "text"}, "config": {"value": "minikube (1)"}},
          "Last Seen": {"metadata": {"type": "timestamp"}, "config": {"timestamp": 1592272167}},
          "Message": {
            "metadata": {
              "type": "link",
              "title": [{"metadata": {"type": "text"}, "config": {"value": ""}}]
            },
            "config": {
              "value": "Node minikube event: Registered Node minikube in Controller",
              "ref": "/overview/namespace/default/events/minikube.1618e3199521f52d"
            }
          },
          "Reason": {"metadata": {"type": "text"}, "config": {"value": "RegisteredNode"}},
          "Type": {"metadata": {"type": "text"}, "config": {"value": "Normal"}}
        }],
        "emptyContent": "We couldn't find any events!",
        "loading": false,
        "filters": {}
      }
    }]
  }
};
