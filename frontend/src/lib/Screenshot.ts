import * as htmlToImage from "html-to-image";
import {
  AutoScreenshot,
  LogError,
  ManualScreenshot,
} from "wailsjs/go/main/App";

export namespace Screenshot {
  export const manual = async (
    targetElementID: string,
    filename: string,
  ): Promise<boolean> => {
    const image = await getBase64Image(targetElementID);
    return await ManualScreenshot(`${filename}.png`, image);
  };

  export const auto = async (
    targetElementID: string,
    filename: string,
  ): Promise<void> => {
    const image = await getBase64Image(targetElementID);
    await AutoScreenshot(`${filename}.png`, image);
  };
}

const getBase64Image = async (targetElementID: string): Promise<string> => {
  try {
    const element = document.getElementById(targetElementID);
    if (!element) {
      throw new Error(`Element with ID ${targetElementID} not found`);
    }

    // Workaround: first screenshot can't draw values in table.
    await htmlToImage.toPng(element);
    return (await htmlToImage.toPng(element)).split(",")[1];
  } catch (error) {
    LogError(`${getBase64Image.name}: ${JSON.stringify(error)}`, {});
    throw error;
  }
};
