import {
  AlgInfo,
  GameType,
  GameProps,
  GameAlias,
  AlgDescr,
  makeRtpList,
  registerGame,
  GameRegistration,
  Gamble
} from '../../../types'; // Main game types
import { NewFlowersGame } from './flowers_rule';
import { calcStatReg } from './flowers_calc';
import { ReelsMap } from './flowers_reels';

const FlowersAliases: GameAlias[] = [
  { prov: "NetEnt", name: "Flowers", date: "2013-11-11" },
];

const FlowersAlgDescr: AlgDescr = {
  gt: GameType.GTslot,
  gp: GameProps.GPlpay |
      GameProps.GPlsel |
      GameProps.GPretrig |
      GameProps.GPfgreel |
      GameProps.GPfgmult |
      GameProps.GPwild,
  sx: 5, // Screen X dimension (reels)
  sy: 3, // Screen Y dimension (rows)
  sn: 17, // Number of symbols in LinePay (seems to be max symbol ID from flowers_rule.go LinePay array size)
  ln: 30, // Number of bet lines (from BetLines.length in flowers_rule.go)
  bn: 0,  // Number of bonus games (0 for Flowers, as per Go)
  rtp: makeRtpList(ReelsMap), // Create RTP list from ReelsMap
};

export const FlowersGameInfo: AlgInfo = {
  aliases: FlowersAliases,
  algDescr: FlowersAlgDescr,
  // gameFactory and calcStat will be part of the registration
};

// Equivalent to Go's init() function for game registration
function registerFlowersGame(): void {
  const registration: GameRegistration = {
    info: FlowersGameInfo,
    factory: NewFlowersGame as () => Gamble, // Ensure NewFlowersGame matches Gamble type
    statCalc: (ctx: any, mrtp: number) => calcStatReg(ctx, mrtp) // Wrap for context if needed
  };
  registerGame("NetEnt/Flowers", registration);
}

// Call the registration function to make the game available
registerFlowersGame();

// Optional: export a function to get the info, or rely on the registration system
export function getFlowersGameInfo(): AlgInfo {
  return FlowersGameInfo;
}
