import {
    Context,
    Reels5x,
    Stat,
    FindClosest,
    ScanReels5x,
    NewStat
} from "../../slot"; // Adjusted path
import { Game, NewGame, ReelsMap } from "./fruitshop_rule";

// Helper type for the writer function used in calc
type StatWriter = { write: (msg: string) => void };

export function CalcStatBon(ctx: Context | null, mrtp: number): number {
    if (Object.keys(ReelsMap).length === 0) {
        console.error("ReelsMap is empty. Cannot calculate bonus stat. Please load reels first.");
        // Attempt to load or initialize ReelsMap here if a default is available, or throw.
        // For now, returning 0 as the Go version might proceed with nil reels leading to 0.
        return 0;
    }
    const { reels } = FindClosest(ReelsMap, mrtp);
    const g = NewGame();
    g.FSR = 5; // Set free spins mode (arbitrary number, matching Go code)
    const s = NewStat();

    const calc = (w: StatWriter): number => {
        const reshuf = s.Count();
        if (reshuf === 0) { // Avoid division by zero if no spins were counted
            w.write("symbols: rtp(sym) = 0.000000%\n");
            w.write("free spins 0, q = 0, sq = 1/(1-q) = 1.000000\n");
            w.write("free games frequency: 1/0\n");
            w.write("RTP = sq*rtp(sym) = 1*0 = 0.000000%\n");
            return 0;
        }

        const costInfo = g.Cost(); // { cost: number, isJp: boolean }
        const cost = costInfo.cost;

        const lrtp = s.LineRTP(cost);
        const srtp = s.ScatRTP(cost);
        const rtpsym = lrtp + srtp;

        const q = s.FreeCount() / reshuf;
        const sq = q >= 1 ? Infinity : 1 / (1 - q); // Handle q >= 1 to avoid negative sq or division by zero in typical case

        let rtp = sq * rtpsym;
        if (sq === Infinity && rtpsym === 0 && q > 0) rtp = 100; // Or some other representation of guaranteed retrigger with no base pay
        else if (sq === Infinity && rtpsym > 0) rtp = Infinity;


        w.write(`symbols: rtp(sym) = ${rtpsym.toFixed(6)}%\n`);
        w.write(`free spins ${s.FreeCount()}, q = ${q.toExponential(5)}, sq = 1/(1-q) = ${sq.toFixed(6)}\n`);
        const freeGamesFreq = s.FreeHits() > 0 ? reshuf / s.FreeHits() : 0;
        w.write(`free games frequency: 1/${freeGamesFreq.toExponential(5)}\n`);
        w.write(`RTP = sq*rtp(sym) = ${sq.toFixed(5)}*${rtpsym.toFixed(5)} = ${rtp.toFixed(6)}%\n`);
        return rtp;
    };

    // Use a basic console writer for the calc function's output
    const consoleWriter: StatWriter = { write: (msg: string) => console.log(msg) };

    // If ScanReels5x is truly async, this needs to be handled.
    // For now, assuming it's synchronous or its async nature is managed internally for this script's purpose.
    return ScanReels5x(ctx, s, g, reels, calc);
}

export function CalcStatReg(ctx: Context | null, mrtp: number): number {
    console.log("*bonus reels calculations*");
    const rtpfs = CalcStatBon(ctx, mrtp);

    if (ctx && ctx.Err && ctx.Err() !== null) {
        return 0;
    }

    console.log("*regular reels calculations*");
    if (Object.keys(ReelsMap).length === 0) {
        console.error("ReelsMap is empty. Cannot calculate regular stat. Please load reels first.");
        return 0;
    }
    const { reels } = FindClosest(ReelsMap, mrtp);
    const g = NewGame();
    const s = NewStat();

    const calc = (w: StatWriter): number => {
        const reshuf = s.Count();
        if (reshuf === 0) {
            w.write("symbols: rtp(sym) = 0.000000%\n");
            w.write("free spins 0, q = 0, sq = 1/(1-q) = 1.000000\n"); // sq is not used directly here but good to show consistency
            w.write("free games frequency: 1/0\n");
            w.write(`RTP = 0.00000(sym) + 0*${rtpfs.toFixed(5)}(fg) = 0.000000%\n`);
            return 0;
        }

        const costInfo = g.Cost();
        const cost = costInfo.cost;
        const lrtp = s.LineRTP(cost);
        const srtp = s.ScatRTP(cost);
        const rtpsym = lrtp + srtp;
        const q = s.FreeCount() / reshuf;
        // const sq = q >= 1 ? Infinity : 1 / (1 - q); // sq is not directly used in final RTP formula for regular spins

        const rtp = rtpsym + q * rtpfs;

        w.write(`symbols: rtp(sym) = ${rtpsym.toFixed(6)}%\n`);
        w.write(`free spins ${s.FreeCount()}, q = ${q.toExponential(5)}\n`); // Removed sq from this line as it's not used in the final RTP calc here
        const freeGamesFreq = s.FreeHits() > 0 ? reshuf / s.FreeHits() : 0;
        w.write(`free games frequency: 1/${freeGamesFreq.toExponential(5)}\n`);
        w.write(`RTP = ${rtpsym.toFixed(5)}(sym) + ${q.toFixed(5)}*${rtpfs.toFixed(5)}(fg) = ${rtp.toFixed(6)}%\n`);
        return rtp;
    };

    const consoleWriter: StatWriter = { write: (msg: string) => console.log(msg) };
    return ScanReels5x(ctx, s, g, reels, calc);
}
