# [pongo](https://en.wikipedia.org/wiki/Pongo_%28genus%29)2

[![GoDoc](https://godoc.org/github.com/iris-contrib/pongo2?status.svg)](https://godoc.org/github.com/iris-contrib/pongo2)
[![Build Status](https://travis-ci.org/iris-contrib/pongo2.svg?branch=master)](https://travis-ci.org/iris-contrib/pongo2)

pongo2 is the django template syntax for [Iris](https://github.com/kataras/iris).

## First impression of a template

```HTML+Django
<html><head><title>Our admins and users</title></head>
{# This is a short example to give you a quick overview of pongo2's syntax. #}

{% macro user_details(user, is_admin=false) %}
	<div class="user_item">
		<!-- Let's indicate a user's good karma -->
		<h2 {% if (user.karma >= 40) || (user.karma > calc_avg_karma(userlist)+5) %}
			class="karma-good"{% endif %}>
			
			<!-- This will call user.String() automatically if available: -->
			{{ user }}
		</h2>

		<!-- Will print a human-readable time duration like "3 weeks ago" -->
		<p>This user registered {{ user.register_date|naturaltime }}.</p>
		
		<!-- Let's allow the users to write down their biography using markdown;
		     we will only show the first 15 words as a preview -->
		<p>The user's biography:</p>
		<p>{{ user.biography|markdown|truncatewords_html:15 }}
			<a href="/user/{{ user.id }}/">read more</a></p>
		
		{% if is_admin %}<p>This user is an admin!</p>{% endif %}
	</div>
{% endmacro %}

<body>
	<!-- Make use of the macro defined above to avoid repetitive HTML code
	     since we want to use the same code for admins AND members -->
	
	<h1>Our admins</h1>
	{% for admin in adminlist %}
		{{ user_details(admin, true) }}
	{% endfor %}
	
	<h1>Our members</h1>
	{% for user in userlist %}
		{{ user_details(user) }}
	{% endfor %}
</body>
</html>
```

# Documentation

For a documentation on how the templating language works you can [head over to the Django documentation](https://docs.djangoproject.com/en/dev/topics/templates/). pongo2 aims to be compatible with it.

You can access pongo2's API documentation on [godoc](https://godoc.org/github.com/iris-contrib/pongo2).

## Caveats 

### Filters

 * **date** / **time**: The `date` and `time` filter are taking the Golang specific time- and date-format (not Django's one) currently. [Take a look on the format here](http://golang.org/pkg/time/#Time.Format).
 * **stringformat**: `stringformat` does **not** take Python's string format syntax as a parameter, instead it takes Go's. Essentially `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`.
 * **escape** / **force_escape**: Unlike Django's behaviour, the `escape`-filter is applied immediately. Therefore there is no need for a `force_escape`-filter yet.

### Tags

 * **for**: All the `forloop` fields (like `forloop.counter`) are written with a capital letter at the beginning. For example, the `counter` can be accessed by `forloop.Counter` and the parentloop by `forloop.Parentloop`.
 * **now**: takes Go's time format (see **date** and **time**-filter).

### Misc

 * **not in-operator**: You can check whether a map/struct/string contains a key/field/substring by using the in-operator (or the negation of it):
    `{% if key in map %}Key is in map{% else %}Key not in map{% endif %}` or `{% if !(key in map) %}Key is NOT in map{% else %}Key is in map{% endif %}`.

# Add-ons, libraries and helpers

## Official

 * [pongo2-addons](https://github.com/iris-contrib/pongo2-addons) - Official additional filters/tags for pongo2 (for example a **markdown**-filter). They are in their own repository because they're relying on 3rd-party-libraries.

# API-usage examples

Please see the documentation for a full list of provided API methods.

Please refer to [kataras/iris/_examples/view](https://github.com/kataras/iris/tree/master/_examples/view)

## Contributors

This project exists thanks to all the people who contribute. 
<a href="https://github.com/iris-contrib/pongo2/graphs/contributors"><img src="https://opencollective.com/pongo2/contributors.svg?width=890&button=false" /></a>
