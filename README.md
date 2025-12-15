# htmangl

A tool for combining HTML files.

This isn't really that smart. It is meant for fairly coarse smashing together of
HTML, if you need precision look elsewhere.

But, for instance, let's say you have a template that you want to apply to some
pages on a website. So you write `template.html`:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>My website</title>
    <link rel="stylesheet" href="styles.css" type="text/css" />
  </head>
  <body>
    <header>
      <h1>My website</h1>
    </header>

    <!-- htmangl:insert -->

    <footer>
      Copyright me (this year)
    </footer>
  </body>
</html>
```

and then write some pages that look a bit like `home_partial.html`:

```html
<!DOCTYPE html>
<html>
  <head>
    <title> - Home</title>
  </head>
  <body>
    <p>Hey</p>
  </body>
</html>
```

to end up with what I want:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8"/>
    <title>My website - Home</title>
    <link rel="stylesheet" href="styles.css" type="text/css"/>
  </head>
  <body>
    <header>
      <h1>My website</h1>
    </header>
    
    <p>Hey</p>

    <footer>
      Copyright me (this year)
    </footer>
  </body>
</html>
```

I just need to run `htmangl template.html home_partial.html > home.html`.
