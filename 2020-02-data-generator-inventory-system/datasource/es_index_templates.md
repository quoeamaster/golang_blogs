
# m_supermarket_inventory
PUT _template/m_supermarket_inventory
{
  "index_patterns": ["m_supermarket_inventory*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "stock_in_cost":{
        "type": "half_float"
      },
      "stock_in_quantity": {
        "type": "integer"
      },
      "stock_in_date": {
        "type": "date"
      },
      "expiry_date": {
        "type": "date"
      },
      
      "product": {
        "properties": {
          "id": {
            "type": "keyword"
          },
          "desc": {
            "type": "text"
          },
          "batch_id": {
            "type": "keyword"
          }
        }
      },
      
      "location": {
        "properties": {
          "id": {
            "type": "keyword"
          },
          "name": {
            "type": "text",
            "fields": {
              "raw": {
                "type": "keyword"
              }
            }
          },
          "post_code": {
            "type": "keyword"
          },
          "coord": {
            "type": "geo_point"
          }
        }
      }
    }
  }
}

# m_supermarket_sales
PUT _template/m_supermarket_sales
{
  "index_patterns": ["m_supermarket_sales*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "date": {
        "type": "date"
      },
      "selling_price": {
        "type": "half_float"
      },
      "quantity": {
        "type": "integer"
      },
      
      "product": {
        "properties": {
          "id": {
            "type": "keyword"
          },
          "desc": {
            "type": "text"
          },
          "batch_id": {
            "type": "keyword"
          }
        }
      },
      
      "client": {
        "properties": {
          "id": {
            "type": "keyword"
          },
          "name": {
            "type": "text",
            "fields": {
              "raw": {
                "type": "keyword"
              }
            }
          },
          "gender": {
            "type": "keyword"
          },
          "occupation": {
            "type": "text",
            "fields": {
              "raw": {
                "type": "keyword"
              }
            }
          }
        }
      },
      
      "location": {
        "properties": {
          "id": {
            "type": "keyword"
          },
          "name": {
            "type": "text",
            "fields": {
              "raw": {
                "type": "keyword"
              }
            }
          },
          "post_code": {
            "type": "keyword"
          },
          "coord": {
            "type": "geo_point"
          }
        }
      }
    }
  }
}