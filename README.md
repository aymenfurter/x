X
==================
X is a Go application that lets you leverage the power of OpenAI's large language models to solve a wide range of tasks, all from the comfort of your terminal.
The key difference between X and existing terminal-based copilot solutions is the role reversal. In X, the human user takes on the role of the copilot, while the AI assumes the pilot seat.


Getting Started
---------------

To start using X, follow these simple steps:

1.  Ensure you have a valid OpenAI API key. You can obtain one by signing up at [OpenAI's website](https://www.openai.com/).

2.  Set the `OPEN_AI_KEY` environment variable by running the following command in your terminal:
    `export OPEN_AI_KEY=your_api_key_here`

3.  Clone this repository and navigate to the project directory.

4.  Build and install the X application:
    `go build -o x
    sudo mv x /usr/local/bin`

Usage
-----

Once you have set up the environment variable and installed the X application, you can start using it to solve tasks. For instance, to build a new Node.js CLI app that allows you to manage pets, simply run the following command:
`x "build a new nodejs cli app allows me to manage pets"`

The AI will go through the task step by step, asking for your confirmation before executing each step. This interactive approach ensures that you stay in control while the AI takes care of the heavy lifting.

Demo Video
----------
Watch this demo video to see X in action, solving the Node.js CLI app task mentioned above:

[![X Demo Video](https://img.youtube.com/vi/VIPRT7NlHC8/0.jpg)](https://www.youtube.com/watch?v=VIPRT7NlHC8)


License
-------
Released under Apache License
