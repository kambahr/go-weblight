# A Go website template with multiple master pages

## Weblight is an alternative way of presenting a website: a one-page look&feel -- with operations of a full-feature website with Go managing the background work.

### Three different styles have been included to demonstrate the use of multiple master pages inside the same website. 

#### Overview

The main body of the is loaded (master-page) once; and  all other contents are retrieved in the background.

#### Browser History
The browser history is written via the window.onpopstate() event (for a full user-experience).

#### Navigation

Within the master page, getContent(), javascript, function gets the content (simulating href of an A tag).
To get to the outside of the master page: the <a> tag is used. The web server will whether to server the master page or a single page independently.

#### Javascript
javascript code will not run in any of the content pages by simply placing them in <script> blocks; but it will rather run in-line upon firing events: e.g. onclick="alert('hello world');"

&lt;script&gt; blocks can only be used inside the maser page.

#### Search Engines

Although target content is written to html tags (the entire page not rendered with each request), search engines can still retrieve your full content as that of a traditional website. And it is still possible to make use of <meta name="robots" content="noindex"> on a master page for excluding target page.
                                
#### Some Usage Example

This type of website, basically, loads faster and may offers a better user-experience that comes close to that of a desktop application. Although this style of delivering content can be used for any environment, the following are some suitable examples.

##### User Experience

- When you just don't want your main content-page flickering.
- When you want music running in the background, while navigating within your website.
- Doc sites; where users go back and fourth, quickly and frequently to find content.
- The target pages have heavy backgrounds; therefore, best suited for loading the main page's frame-work (main body) once - and then update parts of the page as required.

And... with Go managing the background work!
