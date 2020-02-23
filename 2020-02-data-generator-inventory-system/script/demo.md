





# supermarket

GET _cat/indices?h=index&s=index

GET _template
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
            "type": "text",
            "fields": {
              "raw": {"type": "keyword"}
            }
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
            "type": "text",
            "fields": {"raw": {"type": "keyword"}}
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


GET m_supermarket_inventory/_search
{
  "query": {
    "match": {
      "product.desc": "dotcom"
    }
  }
}

# different supermarket's inventory data
GET m_supermarket_inventory/_search?filter_path=hits.total.value,aggregations.NAME.buckets.key,aggregations.NAME.buckets.doc_count,aggregations.NAME.buckets.NAME.hits.hits._source
{
  "size": 0, 
  "aggs": {
    "NAME": {
      "terms": {
        "field": "location.post_code",
        "size": 5
      },
      "aggs": {
        "NAME": {
          "top_hits": {
            "_source": ["product.desc"], 
            "size": 10
          }
        }
      }
    }
  }
}
GET m_supermarket_inventory/_search
{
  "sort": [
    {
      "stock_in_quantity": {
        "order": "desc"
      }
    }
  ]
}
GET m_supermarket_inventory/_search
{
  "size": 0, 
  "aggs": {
    "NAME": {
      "terms": {
        "field": "product.desc.raw",
        "size": 5000
      },
      "aggs": {
        "NAME": {
          "terms": {
            "field": "location.name.raw"
          }
        }
      }
    }
  }
}

# sales information (finally)
GET m_supermarket_sales/_search
{
  "query": {
    "geo_distance": {
      "distance": "2000m",
      "location.coord": {
        "lat": 1.3323346,
        "lon": 103.93805
      }
    }
  }
}

GET m_supermarket_sales/_search
{
  "size": 0,
  "aggs": {
    "NAME": {
      "auto_date_histogram": {
        "field": "date"
      },
      "aggs": {
        "shop": {
          "terms": {
            "field": "location.name.raw",
            "size": 4,
            "min_doc_count": 1
          }
        },
        "max_qty": {
          "max": {
            "field": "quantity"
          }
        }
      }
    }
  }
}










