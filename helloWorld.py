import tkinter as tk


def main():
    # Create the main window
    root = tk.Tk()
    root.title("Hello World App")

    # Set the size of the window
    root.geometry("300x200")

    # Add a label with "Hello World"
    label = tk.Label(root, text="Hello World", font=("Helvetica", 16))
    label.pack(expand=True)

    # Start the Tkinter event loop
    root.mainloop()


if __name__ == "__main__":
    main()
