# Baryon Specification

The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in BCP 14
[RFC2119](https://www.ietf.org/rfc/rfc2119.txt)
[RFC8174](https://www.ietf.org/rfc/rfc8174.txt) when, and only when, they
appear in all capitals, as shown here.

Baryon follows the R documentation standard but provides an augmentation layer
to provide context for [Galaxy's Tools](https://galaxyproject.org).

To find more information about documenting your R functions, please refer to
the R's roxygen2 package documentation:

- https://cran.r-project.org/web/packages/roxygen2/

## Baryon Namespace

The tiniest Baryon Namespace is defined as follows:

```
$B{}
▲ ▲
│ │
│ Instructions
│
Prefix
```

Specific instruction for Baryon MUST live inside the Baryon Namespace.
Each roxygen2 tag MAY contain (at most) one Baryon Namespace.
Subsequent Baryon Namespaces will be ignored. It's NOT RECOMMENDED to have
multiple Baryon Namespaces for a single roxygen2 tag.

Instruction inside a Baryon Namespace MUST be separated by a semicolon `;`.
Whitespace that might occur before or after an instruction is ignored.
The last instruction MAY have a delimiting semicolon.
Baryon Namespaces including only one instruction MAY have a delimiting
semicolon.

A Baryon Instruction MUST be a sequence of alphabetical characters or MUST
be the `!` special character, and, MAY have an argument (or a list of
arguments).

The high level Baryon Specification doesn't specify how arguments must be
subdivided and relegates this definition to the specific implementation of
an instruction. Although, it is RECOMMENDED to have a comma-separated list.

Instruction may override previous instructions and are applied in a
first-come-first-serve fashion.
It is NOT RECOMMENDED to have the same instruction repeated in the same Baryon
Namespace.

## Instructions - Parameters

### required

`required` instruction tags a parameter as required.

Alias: `!`.

Example(s):
```
$B{required}
$B{required;}
$B{!}
```

### type

`type` instructions tags a parameter with its expected type.
Accepts a parameter.

Possible values are: `text`, `integer`, `float`, `boolean`, `genomebuild`,
`select`, `color`, `data_column`, `hidden`, `hidden_data`, `baseurl`, `file`,
`ftpfile`, `data`, `data_collection`, `drill_down`.

Example(s):
```
$B{type(text)}
```

### value

`value` instructions tags a parameter with its default value.
Accepts a parameter.

Example(s):
```
$B{value(a random string)}
```

### options

`options` instructions tags a parameter with its possible options.
Accepts a parameter list.
The last element MAY have a delimiting comma.

Example(s):
```
$B{options(a,random,list,of,parameters)}
$B{options(
    a,
    random,
    list,
    of,
    parameters,
)}
$B{options(
    a,
    random,
    list
)}
$B{options()}
```

## Full example

```
$B{
    type(integer);
    value(4);
    options(
        1,2,3,4,
        5,6,7,8,
    );
    required;
}
```

## Instructions - Description

### container

`container` tags the container that will be used by tool. Accepts three parameters:  
- `<name>` - the container name. Required.
- `<type>` - the containerization technology that will be used. Optional, defaults to `docker`.

Example(s):
```
${container(hello-world:latest,docker)}
${container(hello-world:latest)}
```
