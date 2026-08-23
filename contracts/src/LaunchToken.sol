// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract LaunchToken is ERC20 {
    string public contractURI;

    constructor(
        string memory name_,
        string memory symbol_,
        string memory contractURI_,
        uint256 supply,
        address recipient
    ) ERC20(name_, symbol_) {
        contractURI = contractURI_;
        _mint(recipient, supply);
    }
}
