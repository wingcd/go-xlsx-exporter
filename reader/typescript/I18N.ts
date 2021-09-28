import { DataAccess, DataTable } from "./DataAccess";
import { Language, Language_ARRAY } from "./Lang";

export class I18N extends DataTable
{
    private static _inst;
    private static get inst(): I18N {
        if(!this._inst) {
            this._inst = new I18N(Language_ARRAY);
        }
        return this._inst;
    }

    private constructor(dataType:Function, keyName = "ID") {
        super(dataType, keyName);

        this.initial(this.dataType, null, I18N.getFilename);
    }

    public static currentLanguage: string = "cn";
    public static setLanguage(lan = "cn")
    {
        if (this.currentLanguage != "")
        {
            this.inst.clear();
        }

        this.currentLanguage = lan;
    }

    protected static getFilename(typeName: string): string
    {
        typeName = typeName.replace("_ARRAY", "");
        var datafile = DataAccess.dataDir + typeName.toLocaleLowerCase();
        return `${datafile}.${I18N.currentLanguage}`;
    }

    public static translate(key: string): string
    {
        let lan:Language = this.inst.itemMap[key];
        if (lan)
        {
            return lan.Text;
        }

        return null;
    }
}
