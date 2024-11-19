from PIL import Image, ImageDraw, ImageFont
import io
import base64
import requests
from pdf2image import convert_from_path


def pdf_to_image_preview(file_path):
    images = convert_from_path(file_path, dpi=200, first_page=1, last_page=1)
    if images:
        # Convert the PIL Image to a bytes object in PNG format
        buffered = io.BytesIO()
        images[0].save(buffered, format="PNG")
        img_byte_data = buffered.getvalue()

        # Convert byte data to Base64
        encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

        return encoded_img

    return ""


def text_to_image_preview(file_path, num_chars=1000):
    with open(file_path, "r") as file:
        text = file.read(num_chars)

    # Create a blank image
    image = Image.new("RGB", (256, 256), color=(255, 255, 255))
    d = ImageDraw.Draw(image)
    fnt = ImageFont.truetype(
        "/code/app/Roboto-Regular.ttf", 15
    )  # Specify a font here
    d.text((10, 10), text, font=fnt, fill=(0, 0, 0))

    # Convert the PIL Image to a bytes object in PNG format
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    img_byte_data = buffered.getvalue()

    # Convert byte data to Base64
    encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

    return encoded_img


def youtube_to_image_preview(youtube_url):
    # Extract the video ID from the URL
    video_id = youtube_url.split("v=")[1].split("&")[0]

    # Construct the thumbnail URL
    thumbnail_url = f"https://img.youtube.com/vi/{video_id}/maxresdefault.jpg"

    # Fetch the thumbnail image
    response = requests.get(thumbnail_url)
    image = Image.open(io.BytesIO(response.content))

    # Convert the PIL Image to a bytes object in PNG format
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    img_byte_data = buffered.getvalue()

    # Convert byte data to Base64
    encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

    return encoded_img
