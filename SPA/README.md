# SPA

This SPA folder contains sample snippets for developing production-quality Single-Page Applications in Vuejs.

## Organization
The SPA folder is organized into the following files and folders.

### App.vue
App.vue contains the main template for the SPA. The navigation bar routes to all available views and adapts according to the app's state (i.e., whether or not the user is logged in). For simplicity the App.vue file is the only one containing Bulma classes and custom CSS.

### Assets
The assets folder contains the main.scss file for importing a CSS framework, in this case Bulma, and defining custom variables for the SPA.

### Components
The components folder contains the components for displaying user information and posts.

### Main.js
Main.js defines our main Vue instance with all dependencies. Global settings for Axios, the library used for AJAX API calls, are also defined here (e.g., common headers and instructions for every call).

### Router.js
Router.js defines the Vue Router instance and all routes and route guards. Only routes defined as public can be accessed freely whereas all others require the user to be logged in.

### Store.js
Store.js defines the Vuex instance. Vuex acts as the universal store, allowing user state (i.e., logged in or not) to be preserved and shared across all parts of the SPA.

### Views
The views folder contains the views the router routes to. These files call the API and load the data into imported components.

### Vue.config.js
Vue.config.js sets a couple configurations for development: a reverse proxy for communicating with the API and a source-map for debugging.

## Improvements

This repository would benefit from several improvements.

### Store.js
Currently Store.js has one mutation and one action: updateUser. These same functions are used for both logging in and logging out. To improve the clarity of the code throughout the SPA, a potential improvement would be to separate these into more distinct functions, each serving either to set or remove user state, but not both.

### Error Messages
The error messages are currently quite simple. The benefit is not revealing too much information about the internal code and processes. The messages could nevertheless be more descriptive about potential actions to take, for example prompting users to try again. The display of errors could also be clearer, using certain colours and effects (e.g., borders) to draw attention to the errors.

### Tests
Units tests must be added.

### Accessibility
The UI needs to be accessible to everyone. Guidelines for improving accessibility should be followed, for example always using alternative text on HTML image tags for those who use screen readers.

### Overall Expansion
The improvements suggested above are only a few key issues that stand out with the code presented here. Of course many improvements could be made across the SPA, such as supporting international languages, improving the UX, greatly extending functionality, etc. These snippets are only meant to demonstrate useful design patterns for developing production-quality SPAs: the potential implementations and extensions are limitless.
