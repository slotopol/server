import {
  Reels5x,
  SlotGame,
  Stat,
  Wins,
  Screen5x3, // Assuming Game uses Screen5x3
  Reels,   // For adapting Reels5x
  readObj
} from '../../slot/types';
import { FlowersGame, NewFlowersGame } from './flowers_rule';
import { ReelsBon, ReelsMap } from './flowers_reels';

// Simplified context object for calculations, as Go's context is for cancellation.
interface CalcContext {
  isCancelled: () => boolean;
  cancel?: () => void; // Optional cancel function
}

/**
 * Simulates the core logic of slot.ScanReels5x for statistical calculation.
 * Iterates through all reel combinations for a 5-reel game.
 * @param ctx Calculation context for potential cancellation.
 * @param stat Stat object to update.
 * @param gameInstance A fresh game instance to use for scanning.
 * @param reels The reels to scan.
 * @param calcCallback The callback function to execute for each combination, receiving a writer (here, a stats summary).
 * @returns The final RTP calculated by the callback.
 */
function scanReels5xSimulation(
  ctx: CalcContext,
  stat: Stat,
  gameInstance: FlowersGame, // Use FlowersGame or a more generic SlotGame if applicable
  reelsData: Reels5x,
  calcCallback: (statSummary: string) => number
): number {
  const r1 = reelsData[0];
  const r2 = reelsData[1];
  const r3 = reelsData[2];
  const r4 = reelsData[3];
  const r5 = reelsData[4];

  if (!r1 || !r2 || !r3 || !r4 || !r5) {
    console.error("Reel data is incomplete.");
    return 0;
  }

  const totalReshuffles = r1.length * r2.length * r3.length * r4.length * r5.length;
  stat.setPlan(totalReshuffles);

  // Adapt Reels5x (Sym[][]) to Reels interface for screen.setCol
    const reelsAdapter: Reels = {
        cols: () => 5,
        reel: (colPos) => reelsData[colPos-1], // col is 1-indexed
        reshuffles: () => totalReshuffles,
        toString: () => reelsData.map(r => r.length).join(', ')
    };

  let reshufflesDone = 0;
  const CtxGranulation = 100000; // Check for cancellation less frequently for performance

  for (let i1 = 0; i1 < r1.length; i1++) {
    gameInstance.screen.setCol(1, r1, i1);
    for (let i2 = 0; i2 < r2.length; i2++) {
      gameInstance.screen.setCol(2, r2, i2);
      for (let i3 = 0; i3 < r3.length; i3++) {
        gameInstance.screen.setCol(3, r3, i3);
        for (let i4 = 0; i4 < r4.length; i4++) {
          gameInstance.screen.setCol(4, r4, i4);
          for (let i5 = 0; i5 < r5.length; i5++) {
            reshufflesDone++;
            if (reshufflesDone % CtxGranulation === 0) {
              if (ctx.isCancelled()) {
                console.log("Calculation cancelled.");
                return calcCallback(getStatSummary(stat, gameInstance.cost().cost, 0, 0)); // Return current RTP
              }
            }
            gameInstance.screen.setCol(5, r5, i5);

            const wins: Wins = [];
            gameInstance.scanner(wins); // Populates wins
            stat.update(wins, 1); // cfn = 1 for regular spins in this context
          }
        }
      }
    }
  }
  return calcCallback(getStatSummary(stat, gameInstance.cost().cost, 0, 0));
}

function getStatSummary(stat: Stat, cost: number, qVal?: number, rtpFsVal?: number): string {
  const reshuf = stat.getCount();
  const lrtp = reshuf > 0 ? stat.lineRTP(cost) : 0;
  const srtp = reshuf > 0 ? stat.scatRTP(cost) : 0;
  const rtpsym = lrtp + srtp;

  let q = qVal !== undefined ? qVal : (reshuf > 0 ? stat.getFreeCount() / reshuf : 0);
  let sq = q < 1 ? 1 / (1 - q) : Infinity; // Avoid division by zero if q >= 1

  let rtp: number;
  let summary: string;

  if (rtpFsVal !== undefined) { // Regular game calc
    rtp = rtpsym + q * rtpFsVal;
    summary = `symbols: ${lrtp.toFixed(5)}(lined) + ${srtp.toFixed(5)}(scatter) = ${rtpsym.toFixed(6)}%\n` +
              `free spins ${stat.getFreeCount().toExponential(8)}, q = ${q.toFixed(5)}, sq = 1/(1-q) = ${sq.toFixed(6)}\n` +
              `free games frequency: 1/${(reshuf / stat.getFreeHits() || 0).toFixed(5)}\n` +
              `RTP = ${rtpsym.toFixed(5)}(sym) + ${q.toFixed(5)}*${rtpFsVal.toFixed(5)}(fg) = ${rtp.toFixed(6)}%`;
  } else { // Bonus game calc
    rtp = sq * rtpsym;
     summary = `symbols: ${lrtp.toFixed(5)}(lined) + ${srtp.toFixed(5)}(scatter) = ${rtpsym.toFixed(6)}%\n` +
              `free spins ${stat.getFreeCount().toExponential(8)}, q = ${q.toFixed(5)}, sq = 1/(1-q) = ${sq.toFixed(6)}\n` +
              `free games frequency: 1/${(reshuf / stat.getFreeHits() || 0).toFixed(5)}\n` +
              `RTP = sq*rtp(sym) = ${sq.toFixed(5)}*${rtpsym.toFixed(5)} = ${rtp.toFixed(6)}%`;
  }
  // console.log(summary); // For debugging or direct output
  return summary; // The calc function in Go returns RTP, so we return it here
}


export function calcStatBon(ctx: CalcContext): number {
  const reels = readObj(ReelsBon); // Use readObj, assuming ReelsBon is already loaded
  const g = NewFlowersGame();
  g.sel = 1; // As per Go code
  g.fsr = 10; // Set free spins mode
  const s = new Stat();

  const calc = (statSummaryOutput: string): number => {
    // Parse RTP from summary or recalculate based on Stat object `s`
    // For simplicity, let's directly use `s` to calculate RTP like in getStatSummary
    const cost = g.cost().cost;
    const reshuf = s.getCount();
    if (reshuf === 0) return 0;

    const lrtp = s.lineRTP(cost);
    const srtp = s.scatRTP(cost);
    const rtpsym = lrtp + srtp;
    const q = s.getFreeCount() / reshuf;
    const sq = q < 1 ? 1 / (1 - q) : Infinity;
    const rtp = sq * rtpsym;

    // The Go version prints to a writer. We can log it or return parts of it.
    // console.log(statSummaryOutput); // Output the summary generated by getStatSummary
    return rtp;
  };

  // return scanReels5xSimulation(ctx, s, g, reels, (summary) => calc(summary));
   return scanReels5xSimulation(ctx, s, g, reels, () => {
    const cost = g.cost().cost;
    const summary = getStatSummary(s, cost); // Generate summary string for logging if needed
    // console.log(summary);
    return calc(summary); // Pass summary for consistency, though calc re-evaluates
  });
}

export function calcStatReg(ctx: CalcContext, mrtp: number): number {
  console.log("*bonus reels calculations*");
  const rtpfs = calcStatBon(ctx); // Should be 447.464713 based on Go comments
  if (ctx.isCancelled()) {
    return 0;
  }
  console.log(`*regular reels calculations for target RTP: ${mrtp}*`);

  // FindClosest equivalent:
  const reelKeys = Object.keys(ReelsMap).map(Number).sort((a, b) => Math.abs(a - mrtp) - Math.abs(b - mrtp));
  const bestRtpKey = reelKeys.length > 0 ? reelKeys[0] : null;
  let reels: Reels5x;

  if (bestRtpKey !== null) {
    reels = ReelsMap[bestRtpKey];
    console.log(`Using reels for actual RTP: ${bestRtpKey}`);
  } else {
    console.error("No reels found for any RTP. Cannot proceed with regular calculation.");
    // Fallback to bonus reels if no regular reels are found, though this is not ideal.
    // Or, more appropriately, throw an error or return a specific value indicating failure.
    // For now, let's use the first available reel set from ReelsMap as a default if any exist.
    const firstKey = Object.keys(ReelsMap)[0];
    if (firstKey) {
        console.warn(`Falling back to first available reels: ${firstKey}`);
        reels = ReelsMap[Number(firstKey)];
    } else {
        console.error("No regular reels available in ReelsMap.");
        return 0; // Or throw
    }
  }

  const g = NewFlowersGame();
  g.sel = 1; // As per Go code
  const s = new Stat();

  const calc = (statSummaryOutput: string): number => {
    const cost = g.cost().cost;
    const reshuf = s.getCount();
    if (reshuf === 0) return 0;

    const lrtp = s.lineRTP(cost);
    const srtp = s.scatRTP(cost);
    const rtpsym = lrtp + srtp;
    const q = s.getFreeCount() / reshuf;
    // const sq = 1 / (1 - q); // Not used for regular RTP calc in Go
    const rtp = rtpsym + q * rtpfs;

    // console.log(statSummaryOutput);
    return rtp;
  };

//   return scanReels5xSimulation(ctx, s, g, reels, (summary) => calc(summary));
  return scanReels5xSimulation(ctx, s, g, reels, () => {
    const cost = g.cost().cost;
    const summary = getStatSummary(s, cost, undefined, rtpfs); // Pass rtpfs for regular calc summary
    // console.log(summary);
    return calc(summary);
  });
}
