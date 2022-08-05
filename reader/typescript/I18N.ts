import { DataAccess, DataTable, IDataArray, IMessage } from "./DataAccess";
import { Language, Language_ARRAY } from "./Lang";

export class I18N extends DataTable<Language>
{
    private static _inst;
    static get inst(): I18N {
        if(!this._inst) {
            this._inst = new I18N();
        }
        return this._inst;
    }

    private constructor(keyName = "ID") {
        super(Language, keyName);

        this.initial(this.dataType, null, I18N.getFilename);
    }

    protected load() : Language[] {
        var buffer = this.onLoadData("Language_ARRAY");
        var message = Language_ARRAY.decode(buffer);
        return (message as IDataArray).Items;        
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

    public static translate(key: string|number): string
    {
        let lan:Language = this.inst.itemMap[key];
        if (lan)
        {
            return lan.Text;
        }

        return null;
    }
}
