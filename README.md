# A Go website template with multiple master pages

## Weblight is an alternative way of presenting a website: a one-page look&feel -- with operations of a full-feature website with Go managing the background work.

### Three different styles have been included to demonstrate the use of multiple master pages inside the same website. 

#### Overview

The main body of the page is loaded (master-page) once; and  all other contents are retrieved in the background.

#### Browser History
The browser history is written via the window.onpopstate() event (for a full user-experience).

#### Navigation

Within the master page, getContent(), javascript, function gets the content (simulating href of an A tag); and
to get to the outside of the master page: the &lt;a&gt; tag is used. The web server will know whether to responde with a master page or a single (independent) page.

#### Javascript

javascript code will not run in any of the content pages by simply placing them in <script> blocks; but it will run in-line, upon firing events: e.g. onclick="alert('hello world');"  &lt;script&gt; blocks can only be used inside the maser page.

#### Search Engines

Although content is rendered partially (the entire page not rendered with each request), search engines will still be able to retrieve the full content as that of a traditional website. It is also still possible to make use of &lt;meta name="robots" content="noindex"&gt; on a master page for excluding target pages from *search* - independently for each content-page.

#### Some Usage Example

This type of website, basically, loads faster and may offer a user-experience close to that of a desktop application. Although this style of delivering content can be used for any environment, the following are some suitable examples.

##### User Experience

- When flickering effect is noticeable, while navigating inside the website.
- When you want music running in the background, while navigating within your website.
- Doc sites; where users go back and fourth, quickly and frequently to find content.

Fore more details, run the sample and see *notes* on the /bare path.