import {
    GameAlias,
    AlgDescr,
    AlgInfo,
    GT,
    GP,
    Date,
    MakeRtpList,
    AlgInfoImpl
} from "../../../linkdata"; // Adjusted path
import { NewGame, ReelsMap } from "./fruitshop_rule";
import { CalcStatReg } from "./fruitshop_calc";
import { LinePay, BetLines } from "./fruitshop_rule";


const Aliases: GameAlias[] = [
    { Prov: "NetEnt", Name: "Fruit Shop", Date: Date(2011, 9, 15) },
];

const Descr: AlgDescr = {
    GT: GT.SLOT,
    GP: GP.L_PAY |
        GP.RETRIG |
        GP.FG_MULT |
        GP.WILD |
        GP.W_MULT,
    SX: 5, // screen width
    SY: 3, // screen height
    SN: LinePay.length -1, // Number of symbols (LinePay includes an empty placeholder for wild)
    LN: BetLines.length,   // Number of lines
    BN: 0, // Number of bonuses (fruitshop doesn't have separate bonus games in the typical sense)
    RTP: MakeRtpList(ReelsMap), // RTP list derived from ReelsMap keys
};

// Create the AlgInfo object using the implementation that has the SetupFactory method
export const Info: AlgInfo = new AlgInfoImpl(Aliases, Descr);

// Initialize and register the game.
// This replaces the Go init() function.
// This should be called once when the game module is loaded.
function initializeGame() {
    // Update RTP list in case ReelsMap was populated after Descr was defined.
    // This is important if ReelsMap is loaded dynamically.
    Descr.RTP = MakeRtpList(ReelsMap);
    Info.SetupFactory(NewGame, CalcStatReg);
}

// Call initialization
initializeGame();

// Optional: A function to re-initialize if ReelsMap changes, though not typical.
export function UpdateGameRegistrationWithReels(newReelsMap: any) {
    // This is an advanced scenario. Typically ReelsMap is static or loaded once.
    // If ReelsMap can be dynamically updated and you need to refresh RTPs:
    Descr.RTP = MakeRtpList(newReelsMap);
    // Note: SetupFactory in the current linkdata.ts adds to AlgList and InfoMap.
    // Re-calling it might lead to duplicates if not handled carefully in linkdata.ts.
    // For simplicity, we assume ReelsMap is stable after initial load for registration.
    // If re-registration is needed, linkdata.ts might need a way to update or replace existing entries.
    console.log("Game registration RTPs updated based on new ReelsMap.");
}
