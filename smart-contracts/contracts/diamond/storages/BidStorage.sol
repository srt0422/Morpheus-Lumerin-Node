// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Paginator} from "@solarity/solidity-lib/libs/arrays/Paginator.sol";

import {IBidStorage} from "../../interfaces/storage/IBidStorage.sol";

contract BidStorage is IBidStorage {
    using Paginator for *;
    using EnumerableSet for EnumerableSet.Bytes32Set;

    struct BidsStorage {
        address token; // MOR token
        mapping(bytes32 bidId => Bid) bids; // bidId = keccak256(provider, modelId, nonce)
        mapping(bytes32 modelId => EnumerableSet.Bytes32Set) modelBids; // keccak256(provider, modelId) => all bidIds
        mapping(bytes32 modelId => EnumerableSet.Bytes32Set) modelActiveBids; // modelId => active bidIds
        mapping(address provider => EnumerableSet.Bytes32Set) providerBids; // provider => all bidIds
        mapping(address provider => EnumerableSet.Bytes32Set) providerActiveBids; // provider => active bidIds
        mapping(bytes32 providerModelId => uint256) providerModelNonce; // keccak256(provider, modelId) => last nonce
    }

    bytes32 public constant BIDS_STORAGE_SLOT = keccak256("diamond.standard.bids.storage");

    /** PUBLIC, GETTERS */
    function getBid(bytes32 bidId_) external view returns (Bid memory) {
        return _getBidsStorage().bids[bidId_];
    }

    function getProviderActiveBids(
        address provider_,
        uint256 offset_,
        uint256 limit_
    ) external view returns (bytes32[] memory, uint256) {
        EnumerableSet.Bytes32Set storage providerActiveBids = _getBidsStorage().providerActiveBids[provider_];

        return (providerActiveBids.part(offset_, limit_), providerActiveBids.length());
    }

    function getModelActiveBids(
        bytes32 modelId_,
        uint256 offset_,
        uint256 limit_
    ) external view returns (bytes32[] memory, uint256) {
        EnumerableSet.Bytes32Set storage modelActiveBids = _getBidsStorage().modelActiveBids[modelId_];

        return (modelActiveBids.part(offset_, limit_), modelActiveBids.length());
    }

    function getProviderBids(
        address provider_,
        uint256 offset_,
        uint256 limit_
    ) external view returns (bytes32[] memory, uint256) {
        EnumerableSet.Bytes32Set storage providerBids = _getBidsStorage().providerBids[provider_];

        return (providerBids.part(offset_, limit_), providerBids.length());
    }

    function getModelBids(
        bytes32 modelId_,
        uint256 offset_,
        uint256 limit_
    ) external view returns (bytes32[] memory, uint256) {
        EnumerableSet.Bytes32Set storage modelBids = _getBidsStorage().modelBids[modelId_];

        return (modelBids.part(offset_, limit_), modelBids.length());
    }

    function getToken() external view returns (address) {
        return _getBidsStorage().token;
    }

    function isBidActive(bytes32 bidId_) public view returns (bool) {
        Bid storage bid = _getBidsStorage().bids[bidId_];

        return bid.createdAt != 0 && bid.deletedAt == 0;
    }

    /** INTERNAL */
    function _isModelActiveBidsEmpty(bytes32 modelId) internal view returns (bool) {
        return _getBidsStorage().modelActiveBids[modelId].length() == 0;
    }

    function _isProviderActiveBidsEmpty(address provider) internal view returns (bool) {
        return _getBidsStorage().providerActiveBids[provider].length() == 0;
    }

    function _getBidsStorage() internal pure returns (BidsStorage storage ds) {
        bytes32 slot_ = BIDS_STORAGE_SLOT;

        assembly {
            ds.slot := slot_
        }
    }
}
