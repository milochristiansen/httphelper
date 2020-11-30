
# HTTP Helper: Make servers easy again!

HTTP Helper is a simple library to make it easier to create simple web servers that need to provide a mix of
static, template, and generated content. This is not suitable for large websites, as all content is loaded and
served from memory.

When the server is initialized all data files are loaded and classified. Once classification is done and all
handlers are assigned their requested resources, anything left over is served as static content.
