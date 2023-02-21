# imgctl

A command-line interface to control images with ease.

## 🚀 Installation

You can install `imgctl` via go:

```bash
go install github.com/cosmtrek/imgctl@latest
```

## 🔧 Usage

To see the available commands, run `imgctl --help`.


### 📜 Extracting text from an image

To extract text from an image using the default OCR service (Tencent), run:

```bash
export TENCENT_OCR_SECRET_ID=your_secret_id
export TENCENT_OCR_SECRET_KEY=your_secret_key

imgctl ocr
```

If the image file path is not set, the program will read the image from the **clipboard**.

To use a specific image file, run:

```bash
imgctl ocr --image /path/to/image.jpg
```

For a full list of available options for the ocr command, run `imgctl ocr --help`.

## 📝 License

This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.