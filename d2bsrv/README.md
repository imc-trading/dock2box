# Filtering

## Query

To make queries just add the field and value, you can specify one or several fields.

> At the moment it only support's equal to.

**Example:**

```
/v1/tags?tag=latest&imageId=568d2ba85d099040397ae363
```

## fields

This allows you to specify which fields you want in the result.

**Example:**
```
/v1/tags?fields=id,tag
```

## sort

This allows you to sort the result ascending or descending, you can specify one or several fields.

**Example ascending:**
```
/v1/tags?sort=tag,created
```

You can sort descending by adding a minus sign.

**Example descending:**
```
/v1/tags?sort=tag,-created
```

# Options

## envelope

This enabled/disables embedding data in an envelope with additional info that normally is available as HTTP status code.

**Example:**
```
/v1/tags?envelope=true
```

**Example output:**
```
{
  "code": 200,
  "data": [
    ...
  ]
}
```

> This can be enabled globally by using **-enable-envelope** when starting the server.

## hateoas

This enables/disables HATEOAS which includes links to methods and related endpoints.

**Example:**
```
/v1/tags?hateoas=true
```

**Example output:**
```
"links": [
  {
    "href": "http://yggdrasil.trading.imc.intra:8080/v1/images/568d2ba85d099040397ae363",
    "rel": "self",
    "method": "GET"
  }
]
```

> This can be disabled globally by using *-disable-hateoas* when starting the server.

## embed

This enables/disables embedding related data in the result. This will affect perfomance since it has the server has to do additional queries.

**Example:**
```
/v1/tags?embed=true
```

**Example output:**
```
{
  "id": "568d2ba85d099040397ae365",
  "tag": "untested",
  "created": "2015-12-01T13:01:05Z",
  "sha256": "37ff8e2ae04a1570781a63a247fce789352beae2889f1d720b2efbec50ef8e0d",
  "imageId": "568d2ba85d099040397ae363",
  "image": {
    "id": "568d2ba85d099040397ae363",
    "image": "test2",
    "type": "docker",
    "bootTagId": "568d2ba85d099040397ae362"
  }
}
```

# Data model

## Hosts

Field | Type | Description
--- | --- | ---
Id | Unique id | Unique Id for host
host | String | Hostname
build | Boolean | If host should be provisioned when PXE booting
debug | Boolean | Debug info and no reboot after provisioning finishes
gpt | Boolean | Use GUID Partition Table
tagId | Foreign key | Foreign key to image/tag
kOpts | String | Kernel options
tenantId | Foreign key | Foreign key to tenant
labels | List of strings | Labels for host
siteId | Foreign key | Foreign key to site

# ROADMAP

- Add checks for foreign keys when creating/updating an entry
- Add more complex queries
- Refactor boiler plate code in each controller and normalize it
- Rename upper-case ID to Id
