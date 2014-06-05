// version: 1.1.0

// 获取各种的随机数.
// 为了加快速度, 用hash算法来产生伪随机数, hash 的 salt 是用 crypto/rand 产生的相对真的随机数.
// 安全着想 salt 定期变化, 当前时间也作为 hash 的一个因子.
package random
