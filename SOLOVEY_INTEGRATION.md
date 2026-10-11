# Narrow pinned session-state synchronization

Semantic source: original anytls/sing-anytls v0.0.11, commit
130d2e61b8895727bfed4942c535e91b246a9603, GPL-3.0-or-later (LICENSE
and upstream authorship retained). This maintained integration is not
byte-identical to the original module. It retains the original module path,
Go1.23.1 declaration, dependency versions, protocol, constructors, authentication,
padding, pool policy and close-once/lifecycle behavior.

Two owner-local fields have concurrent read/write paths: peerVersion is
written by the wire receiver and read by stream creation and handshake
reporting; dieErr is written by stream close and read by concurrent net.Conn
Read/Write/Close. Existing official sing-box v1.14.2 AnyTLS outbound uses
this session client; these are reachable dependency operations.

Use the module's already-required sing/common/atomic types for those fields.
Peer version retains the original byte truncation before storage. Terminal
error retains the original concrete error object and first-close result.
No new dependency, service, registry, credential owner, unsafe/reflection
access, protocol field, retry, source-address identity or privilege is added.
The official core's existing metadata adapter is unchanged and continues to
derive its field offsets from actual compiled types.

The focused session-owner regressions exercise actual server-settings parsing
concurrently with OpenStream, then concurrent stream Close/Write. Unmodified
v0.0.11 fails both tests under supported GCC race detection. The correction
must pass those regressions, Solovey's authenticated parent admission fixture,
and all current required Linux/Windows checks before product integration.

Later original v0.0.12/v0.0.13 includes unrelated directory, constructor,
reuse and frame-write changes and does not synchronize these fields. It is
not automatically substituted for this pinned baseline.
