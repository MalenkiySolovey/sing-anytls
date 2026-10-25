# anytls accepted source copy baseline

Authority: https://github.com/anytls/sing-anytls. Copy: https://github.com/MalenkiySolovey/sing-anytls, role mirror/cache/development. Preserve original history, authorship, license and notices. This is a dependency source copy, not a third core provider.

Selected base: 130d2e61b8895727bfed4942c535e91b246a9603, tree 4e18e8c0f95c4cea2249b4cfacc7a83d91a31f50. Pristine ref: mirror/v0.0.11. Canonical patch branch: solovey/anytls-patches. Exactly one previously accepted source commit 8f440972b353f7ace3bbfbac9cf291938de5ab19, tree 7097f65cc5f010e513de42115cb42b9ef8b6be6a, is retained unchanged. Repository documentation and CI commits change no compiled source. Panel module version/checksums remain unchanged.

Classification: ENGINE_PATCH_REQUIRED. Problem: Negotiated session version and terminal stream errors cross concurrent read/close paths without synchronization.

Panel/provider adaptation is insufficient: Panel can request session retirement, but cannot synchronize private library state or rewrite the stream protocol through its provider adapter.

Affected upstream behavior: session negotiation and concurrent terminal stream state.

Regression: ./session race tests on Linux and Windows, three repetitions. The library owns these tests; Panel separately validates runtime consumption.

Selected upstream base: 130d2e61b8895727bfed4942c535e91b246a9603. Contribution status: upstream submission candidate, prepared but not submitted. Removal condition: accept a reviewed authority base that provides equivalent behavior and passes the same regression without this patch.

The original accepted commit predates the new trailer discipline; its exact ID is explicitly grandfathered by the queue verifier so existing Panel build pins need no rewrite. Future source changes require complete problem/boundary/regression/base/contribution/removal trailers and one logical correction per commit. Generate patch files only for review, contribution, or disclosure. Canonical state is this linear Git queue.

The pristine ref accepts reviewed fast-forward upstream synchronization only. Normal development PRs target solovey patch/topic branches; never merge upstream repeatedly into the patch queue. No historical ref is deleted or rewritten.
