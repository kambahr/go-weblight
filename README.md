# A Go website template with multiple master pages

## A one-page look & feel -- with operations of a full-feature website and Go managing the background work

### Multiple master pages within the same website 

#### Overview

The main body (master-page) is loaded once and all other contents are retrieved in the background.

#### Browser History

The browser history is written via the window.onpopstate() event (for a full user-experience).

#### Navigation

##### Within A Master Page

There are two ways to navigate from (within) a master page.

- Via getContent() - a javascript function that gets the content in the background (simulating href of an A tag).
- Using the href of a &lt;a&gt; tag, which moves away from the current master page (traditional way).
 
The web server is able to distinguish between a master page and a single (independent) page in order to respond accordingly.

#### Javascript

javascript code will not run in any of the content (partial) pages by simply placing them in &lt;script&gt; blocks; but it will run in-line, upon firing events: e.g. onclick="alert('hello world');"  &lt;script&gt; blocks can only be used inside a maser page.

#### Search Engines

Although content is rendered partially (the entire page is not rendered with each request), search engines will still be able to retrieve the full content as that of a traditional website. It is also still possible to make use of &lt;meta name="robots" content="noindex"&gt; on a master page for excluding content pages (partial html) from being *searched* (independently) for each content-page.

#### Some Usage Example

This type of website may offer a user-experience close to that of a desktop application. Although this style of delivering content can be used for any environment, the following are some suitable scenarios.

- When you want music running in the background, while navigating within your website.
- Doc sites; where users navigate back and fourth, quickly and frequently.
- When flickering effect is noticeable.

Fore more details, see notes in master/bare.html; or navigate to /bare and expand *notes*, while the website is running.

#### Running the website

- Start a shell window.
- go build -o weblight && ./weblight --portno 8080
- Navigate to http:&#47;&#47;localhost:8080
