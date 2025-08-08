import girlPointing from "src/assets/images/girl_pointing.png";
import girlSorry from "src/assets/images/girl_sorry.png";
import girlStandby from "src/assets/images/girl_standby.png";

export class Tonako {
  private constructor(private readonly imgPath: string) {}

  static readonly Pointing = new Tonako(girlPointing);
  static readonly Standby = new Tonako(girlStandby);
  static readonly Sorry = new Tonako(girlSorry);

  getImgPath(): string {
    return this.imgPath;
  }
}
