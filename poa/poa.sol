pragma solidity ^0.5.11;

    /// @title Proof of Authority Whitelist Proof of Concept
    /// @author Jon Knight
    /// @author Mosaic Networks
    /// @notice Copyright Mosaic Networks 2019, released under the MIT license
    contract POA_Genesis {

    /// @notice Event emitted when the vote was reached a decision
    /// @param _nominee The address of the nominee
    /// @param _yesVotes The total number of yes votes cast for the nominee to date
    /// @param _noVotes The total number of no votes cast for the nominee to date
    /// @param _accepted The decision, true for added to the whitelist, false for rejected
        event NomineeDecision(
            address indexed _nominee,
            uint  _yesVotes,
            uint _noVotes,
            bool indexed _accepted
        );

    /// @notice Event emitted when a nominee vote is cast
    /// @param _nominee The address of the nominee
    /// @param _voter The address of the person who cast the vote
    /// @param _yesVotes The total number of yes votes cast for the nominee to date
    /// @param _noVotes The total number of no votes cast for the nominee to date
    /// @param _accepted The vote, true for accept, false for rejected
        event NomineeVoteCast(
            address indexed _nominee,
            address indexed _voter,
            uint  _yesVotes,
            uint _noVotes,
            bool indexed _accepted
        );

    /// @notice Event emitted when a nominee is proposed
    /// @param _nominee The address of the nominee
    /// @param _proposer The address of the person who proposed the nominee
        event NomineeProposed(
            address indexed _nominee,
            address indexed _proposer
        );

    /// @notice Event emitted when the eviction vote reached a decision
    /// @param _nominee The address of the nominee
    /// @param _yesVotes The total number of yes votes cast for the nominee to date
    /// @param _noVotes The total number of no votes cast for the nominee to date
    /// @param _accepted The decision, true for eviction, false for rejected eviction
        event EvictionDecision(
            address indexed _nominee,
            uint  _yesVotes,
            uint _noVotes,
            bool indexed _accepted
        );

    /// @notice Event emitted when a eviction vote is cast
    /// @param _nominee The address of the nominee
    /// @param _voter The address of the person who cast the vote
    /// @param _yesVotes The total number of yes votes cast for the nominee to date
    /// @param _noVotes The total number of no votes cast for the nominee to date
    /// @param _accepted The vote, true for evict, false for remain
        event EvictionVoteCast(
            address indexed _nominee,
            address indexed _voter,
            uint  _yesVotes,
            uint _noVotes,
            bool indexed _accepted
        );

    /// @notice Event emitted when a nominee is proposed
    /// @param _nominee The address of the nominee
    /// @param _proposer The address of the person who proposed the nominee
        event EvictionProposed(
            address indexed _nominee,
            address indexed _proposer
        );


    /// @notice Event emitted to announce a moniker
    /// @param _address The address of the user
    /// @param _moniker The moniker of the user
        event MonikerAnnounce(
            address indexed _address,
            bytes32 indexed _moniker
        );


    /// @notice Event emitted on error
    /// @param _address The address of the user
    /// @param _message The error message
        event POA_Error(
            address indexed _address,
            string _message
        );


        struct WhitelistPerson {
          address person;
          uint  flags;
        }

        struct NomineeVote {
          address voter;
          bool  accept;
        }

        struct NomineeElection{
          address nominee;
          address proposer;
          uint yesVotes;
          uint noVotes;
          mapping (address => NomineeVote) vote;
          address[] yesArray;
          address[] noArray;
        }

        mapping (address => WhitelistPerson) whiteList;
        uint whiteListCount;
        address[] whiteListArray;
        mapping (address => NomineeElection) nomineeList;
        address[] nomineeArray;
        mapping (address => bytes32) monikerList;
        mapping (address => NomineeElection) evictionList;
        address[] evictionArray;


/// @notice This is no longer required, but an empty function prevents older monetcli versions with the poa init command erroring
function init () public payable checkAuthorisedModifier(msg.sender)
{

}


/// @notice Modifier to check if a sender is on the white list.
modifier checkAuthorisedModifier(address _address)
{
     require(isWhitelisted(_address), "sender is not authorised");
     _;
}


/// @notice Function exposed for Babble Join authority
function checkAuthorised(address _address) public view returns (bool)
{  // needs check on whitelist to allow original validators to be booted.

    return isWhitelisted(_address);
}

/// @notice Function exposed for Babble Join authority wraps checkAuthorised
function checkAuthorisedPublicKey(bytes32  _publicKey) public view returns (bool)
{
   return checkAuthorised(address(uint160(uint256(keccak256(abi.encodePacked(_publicKey))))));

//    This version works in Solidity 0.4.x, but the extra intermediate steps are required in 0.5.x
//      return checkAuthorised(address(keccak256(abi.encodePacked(_publicKey))));
}

/// @notice wrapper function to check if an address is on the nominee list
/// @param _address the address to be checked
/// @return a boolean value, indicating if _address is on the nominee list
function isNominee(address _address) private view returns (bool)
{
     return (nomineeList[_address].nominee != address(0));
}


/// @notice wrapper function to check if an address is on the eviction list
/// @param _address the address to be checked
/// @return a boolean value, indicating if _address is on the eviction list
function isEvictee(address _address) private view returns (bool)
{
     return (evictionList[_address].nominee != address(0));
}



/// @notice wrapper function to check if an address is on the white list
/// @param _address the address to be checked
/// @return a boolean value, indicating if _address is on the white list
function isWhitelisted(address _address) private view returns (bool)
{
     return (whiteList[_address].person != address(0));
}


 /// @notice private function to add user directly to the whitelist. Used to process the Genesis Whitelist.
 function addToWhitelist(address _address, bytes32 _moniker) private {

     if (! isWhitelisted(_address))   // prevent duplicate whitelist entries
     {
        whiteList[_address] = WhitelistPerson(_address, 0);
        whiteListCount++;
        whiteListArray.push(_address);
        monikerList[_address] = _moniker;
        emit MonikerAnnounce(_address,_moniker);
        emit NomineeDecision(_address, 0, 0, true);  // zero vote counts because there was no vote
     }
     else
     {
         emit POA_Error(_address, "Dup Whitelist");
     }
 }


/// @notice Add a new entry to the nominee list
/// @param _nomineeAddress the address of the nominee
/// @param _moniker the moniker of the new nominee as displayed during the voting process
 function submitNominee (address _nomineeAddress, bytes32 _moniker) public payable // checkAuthorisedModifier(msg.sender)
 {
     if ((! isWhitelisted(_nomineeAddress)) && (! isNominee(_nomineeAddress)) )
     {
        nomineeList[_nomineeAddress] = NomineeElection({nominee: _nomineeAddress, proposer: msg.sender,
                    yesVotes: 0, noVotes: 0, yesArray: new address[](0),noArray: new address[](0) });
        nomineeArray.push(_nomineeAddress);
        monikerList[_nomineeAddress] = _moniker;
        emit NomineeProposed(_nomineeAddress,  msg.sender);
        emit MonikerAnnounce(_nomineeAddress, _moniker);
     }  else {
        if (isWhitelisted(_nomineeAddress)) {
            emit POA_Error(_nomineeAddress, "On Whitelist");
        } else {
            emit POA_Error(_nomineeAddress, "On Nominee list");
        }
     }
 }


/// @notice Add a new entry to the eviction list
/// @param _nomineeAddress the address of the evictee
 function submitEviction (address _nomineeAddress) public payable  checkAuthorisedModifier(msg.sender)
 {

     if ((isWhitelisted(_nomineeAddress)) && (! isEvictee(_nomineeAddress)) )
    {
            evictionList[_nomineeAddress] = NomineeElection({nominee: _nomineeAddress, proposer: msg.sender,
                        yesVotes: 0, noVotes: 0, yesArray: new address[](0),noArray: new address[](0) });
            evictionArray.push(_nomineeAddress);
        //        monikerList[_nomineeAddress] = _moniker;
            emit EvictionProposed(_nomineeAddress,  msg.sender);
        //        emit MonikerAnnounce(_nomineeAddress, _moniker);
    } else {
        if (!isWhitelisted(_nomineeAddress)) {
            emit POA_Error(_nomineeAddress, "Not on Whitelist");
        } else {
            emit POA_Error(_nomineeAddress, "On Evictee list");
        }
     }

 }




 ///@notice Cast a vote for a nominator. Can only be run by people on the whitelist.
 ///@param _nomineeAddress The address of the nominee
 ///@param _accepted Whether the vote is to accept (true) or reject (false) them.
 ///@return returns true if the vote has reached a decision, false if not
 ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.
 function castNomineeVote(address _nomineeAddress, bool _accepted)
        public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult){

     decided = false;
     voteresult = false;

//      Check if open nominee, other checks redundant
     if (isNominee(_nomineeAddress)) {


//      Check that this sender has not voted before. Initial config is no redos - so just reject
         if (nomineeList[_nomineeAddress].vote[msg.sender].voter == address(0)) {
             // Vote is valid. So lets cast the Vote
             nomineeList[_nomineeAddress].vote[msg.sender] = NomineeVote({voter: msg.sender, accept: _accepted });

             // Amend Totals
             if (_accepted)
             {
                 nomineeList[_nomineeAddress].yesVotes++;
                 nomineeList[_nomineeAddress].yesArray.push(msg.sender);
             } else {
                 nomineeList[_nomineeAddress].noVotes++;
                 nomineeList[_nomineeAddress].noArray.push(msg.sender);
             }

             emit NomineeVoteCast(_nomineeAddress, msg.sender,nomineeList[_nomineeAddress].yesVotes,
                   nomineeList[_nomineeAddress].noVotes, _accepted);

             // Check to see if enough votes have been cast for a decision
             (decided, voteresult) = checkForNomineeVoteDecision(_nomineeAddress);
         }
     }
     else
     {   // Not a nominee, so set decided to true
         emit POA_Error(_nomineeAddress,"Not nominee");
         decided = true;
     }

     // If decided, check if on whitelist
     if (decided) {
         voteresult = isWhitelisted(_nomineeAddress);
     }
     return (decided, voteresult);
 }


 ///@notice Cast a vote for an eviction. Can only be run by people on the whitelist.
 ///@param _nomineeAddress The address of the potential evictee
 ///@param _accepted Whether the vote is to evict (true) or remain (false) them.
 ///@return returns true if the vote has reached a decision, false if not
 ///@return only meaningful if the other return value is true, returns true if the nominee is now evicted. false otherwise.
 function castEvictionVote(address _nomineeAddress, bool _accepted)
        public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult){

     decided = false;
     voteresult = false;

//      Check if open nominee, other checks redundant
     if (isEvictee(_nomineeAddress)) {

//      Check that this sender has not voted before. Initial config is no redos - so just reject
         if (evictionList[_nomineeAddress].vote[msg.sender].voter == address(0)) {
             // Vote is valid. So lets cast the Vote
             evictionList[_nomineeAddress].vote[msg.sender] = NomineeVote({voter: msg.sender, accept: _accepted });

             // Amend Totals
             if (_accepted)
             {
                 evictionList[_nomineeAddress].yesVotes++;
                 evictionList[_nomineeAddress].yesArray.push(msg.sender);
             } else {
                 evictionList[_nomineeAddress].noVotes++;
                 evictionList[_nomineeAddress].noArray.push(msg.sender);
             }

             emit EvictionVoteCast(_nomineeAddress, msg.sender,evictionList[_nomineeAddress].yesVotes,
                   nomineeList[_nomineeAddress].noVotes, _accepted);

             // Check to see if enough votes have been cast for a decision
             (decided, voteresult) = checkForEvictionVoteDecision(_nomineeAddress);
         }
     }
     else
     {   // Not a nominee, so set decided to true
         emit POA_Error(_nomineeAddress,"Not nominee");
         decided = true;
     }

     // If decided, check if on whitelist
     if (decided) {
         voteresult = ! isWhitelisted(_nomineeAddress);
     }
     return (decided, voteresult);
 }



 ///@notice This function encapsulates the logic for determining if there are enough votes for a definitive decision
 ///@param _nomineeAddress The address of the NomineeElection
 ///@return returns true if the vote has reached a decision, false if not
 ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.

 function checkForNomineeVoteDecision(address _nomineeAddress) private returns (bool decided, bool voteresult)
 {
     NomineeElection memory election = nomineeList[_nomineeAddress];
     decided = false;
     voteresult = false;


     if (election.noVotes > 0)  // Someone Voted No
     {
         declineNominee(election.nominee);
         decided = true;
         voteresult = false;
     }
     else
     {
         // Requires unanimous approval
         if(election.yesVotes >= whiteListCount)
         {
             acceptNominee(election.nominee);
             decided = true;
             voteresult = true;
         }
     }

     if (decided)
     {
         emit NomineeDecision(_nomineeAddress, election.yesVotes, election.noVotes, voteresult);
     }
     return (decided, voteresult);
 }

 ///@notice This function encapsulates the logic for determining if there are enough votes for a definitive decision
 ///@param _nomineeAddress The address of the EvictionElection
 ///@return returns true if the vote has reached a decision, false if not
 ///@return only meaningful if the other return value is true, returns true if the nominee is not now on the whitelist. false otherwise.

 function checkForEvictionVoteDecision(address _nomineeAddress) private returns (bool decided, bool voteresult)
 {
     NomineeElection memory election = evictionList[_nomineeAddress];
     decided = false;
     voteresult = false;


     if (election.noVotes > 0)  // Someone Voted No
     {
         declineEviction(election.nominee);
         decided = true;
         voteresult = false;
     }
     else
     {
         // Requires unanimous approval
         if(election.yesVotes >= (whiteListCount - 1 ))
         {
             acceptEviction(election.nominee);
             decided = true;
             voteresult = true;
         }
     }

     if (decided)
     {
         emit EvictionDecision(_nomineeAddress, election.yesVotes, election.noVotes, voteresult);
     }
     return (decided, voteresult);
 }


 ///@notice This private function adds the accepted nominee to the whitelist.
 ///@param _nomineeAddress The address of the nominee being added to the whitelist
 function acceptNominee(address _nomineeAddress) private
 {
     if (! isWhitelisted(_nomineeAddress))  // avoid re-adding and corrupting the whiteListCount
     {
       whiteList[_nomineeAddress] = WhitelistPerson(_nomineeAddress, 0);
       whiteListArray.push(_nomineeAddress);
       whiteListCount++;
     } else {
        emit POA_Error(_nomineeAddress,"Already Whitelisted");
     }

 // Remove from nominee list
    removeNominee(_nomineeAddress);
 }

 ///@notice This private function removes the accepted evictee from the whitelist.
 ///@param _nomineeAddress The address of the nominee being added to the whitelist
 function acceptEviction(address _nomineeAddress) private
 {
    deWhiteList(_nomineeAddress);
 // Remove from nominee list
    removeEviction(_nomineeAddress);
 }



 ///@notice This private function adds the removes a user from the whitelist. Not currently used.
 ///@param _address The address of the nominee being removed to the whitelist
 function deWhiteList(address _address) private
 {
     if(isWhitelisted(_address))
     {
         delete(whiteList[_address]);
         whiteListCount--;

         for (uint i = 0; i<whiteListArray.length; i++) {
             if (whiteListArray[i] == _address)
             {  // Replace item to be removed with the last item. Then remove last item.
                 whiteListArray[i] = whiteListArray[whiteListArray.length - 1];
                 delete whiteListArray[whiteListArray.length - 1];
                 whiteListArray.length--;
             break;
             }
         }
     }
 }


// Decline nominee from the nomineeList

 ///@notice This private function removes the declined nominee from the nominee list.
 ///@param _nomineeAddress The address of the nominee being removed from the nominee list
 function declineNominee(address _nomineeAddress) private
 {
      removeNominee(_nomineeAddress);
 }


 ///@notice This private function removes the declined nominee from the nominee list.
 ///@param _nomineeAddress The address of the nominee being removed from the nominee list
 function removeNominee(address _nomineeAddress) private
 {
     removeNomineeVote(_nomineeAddress);
// Remove from Mapping
     delete(nomineeList[_nomineeAddress]);

         for (uint i = 0; i<nomineeArray.length; i++) {
             if (nomineeArray[i] == _nomineeAddress)
             {  // Replace item to be removed with the last item. Then remove last item.
                 nomineeArray[i] = nomineeArray[nomineeArray.length - 1];
                 delete nomineeArray[nomineeArray.length - 1];
                 nomineeArray.length--;
               break;
             }
         }
 }


 ///@notice This private function clears the vote mapping in a Nominee Election record.
 ///@param _nomineeAddress The nomineeElection record to be cleansed
function removeNomineeVote(address _nomineeAddress) private
{
    // Iterate through yes and no array and set the addresses in the
    // records held in the vote mapping to zero address
    // THis means it is treated as unset by out validity checks.
    for (uint j = 0 ; j < nomineeList[_nomineeAddress].yesArray.length; j++) {
        nomineeList[_nomineeAddress].vote[nomineeList[_nomineeAddress].yesArray[j]].voter = address(0);
    }
    for (uint j = 0 ; j < nomineeList[_nomineeAddress].noArray.length; j++) {
        nomineeList[_nomineeAddress].vote[nomineeList[_nomineeAddress].noArray[j]].voter = address(0);
    }


}


 ///@notice This private function clears the vote mapping in a Eviction Election record.
 ///@param _nomineeAddress The nomineeElection record to be cleansed
function removeEvicteeVote(address _nomineeAddress) private
{
    // Iterate through yes and no array and set the addresses in the
    // records held in the vote mapping to zero address
    // THis means it is treated as unset by out validity checks.
    for (uint j = 0 ; j < evictionList[_nomineeAddress].yesArray.length; j++) {
        evictionList[_nomineeAddress].vote[evictionList[_nomineeAddress].yesArray[j]].voter = address(0);
    }
    for (uint j = 0 ; j < evictionList[_nomineeAddress].noArray.length; j++) {
        evictionList[_nomineeAddress].vote[evictionList[_nomineeAddress].noArray[j]].voter = address(0);
    }


}


// Deline evictee from the evictionList

 ///@notice This private function removes the declined nominee from the nominee list.
 ///@param _nomineeAddress The address of the nominee being removed from the nominee list
 function declineEviction(address _nomineeAddress) private
 {
      removeEviction(_nomineeAddress);
 }


 ///@notice This private function removes the declined evictee from the Eviction list.
 ///@param _nomineeAddress The address of the nominee being removed from the eviction list
 function removeEviction(address _nomineeAddress) private
 {

     removeEvicteeVote(_nomineeAddress);

// Remove from Mapping
     delete(evictionList[_nomineeAddress]);

         for (uint i = 0; i<evictionArray.length; i++) {
             if (evictionArray[i] == _nomineeAddress)
             {  // Replace item to be removed with the last item. Then remove last item.
                 evictionArray[i] = evictionArray[evictionArray.length - 1];
                 delete evictionArray[evictionArray.length - 1];
                 evictionArray.length--;
               break;
             }
         }
 }



// Information Section.

 function getNomineeElection(address _address) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
 {
     return (nomineeList[_address].nominee, nomineeList[_address].proposer, nomineeList[_address].yesVotes, nomineeList[_address].noVotes);
 }


 function getEvictionElection(address _address) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
 {
     return (evictionList[_address].nominee, evictionList[_address].proposer, evictionList[_address].yesVotes, evictionList[_address].noVotes);
 }

// Array Section. Functions to support Arrays.

 function getNomineeCount() public view returns (uint count)
 {
     return (nomineeArray.length);
 }

 function getEvictionCount() public view returns (uint count)
 {
     return (evictionArray.length);
 }




function getNomineeAddressFromIdx(uint idx) public view returns (address NomineeAddress)
 {
     require (idx < nomineeArray.length, "Requested address is out of range.");
     return (nomineeArray[idx]);
 }


 function getNomineeElectionFromIdx(uint idx) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
 {
     return (getNomineeElection(getNomineeAddressFromIdx(idx))) ;

 }



 function getEvictionAddressFromIdx(uint idx) public view returns (address NomineeAddress)
 {
     require (idx < evictionArray.length, "Requested address is out of range.");
     return (evictionArray[idx]);
 }


 function getEvictionElectionFromIdx(uint idx) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
 {
     return (getEvictionElection(getEvictionAddressFromIdx(idx))) ;

 }


 function getWhiteListCount() public view returns (uint count)
 {
     return (whiteListArray.length);
 }


 function getWhiteListAddressFromIdx(uint idx) public view returns (address WhiteListAddress)
 {
     require (idx < whiteListArray.length, "Requested address is out of range.");
     return (whiteListArray[idx]);
 }


 function getYesVoteCount(address _nomineeAddress)  public view returns (uint count)
 {
     return (nomineeList[_nomineeAddress].yesArray.length);
 }

 function getYesVoterFromIdx(address _nomineeAddress, uint _idx)  public view returns (address voter)
 {
     require (_idx < nomineeList[_nomineeAddress].yesArray.length, "Requested address is out of range.");
     return (nomineeList[_nomineeAddress].yesArray[_idx]);
 }


 function getNoVoteCount(address _nomineeAddress)  public view returns (uint count)
 {
     return (nomineeList[_nomineeAddress].noArray.length);
 }

 function getNoVoterFromIdx(address _nomineeAddress, uint _idx) public view returns (address voter)
 {
     require (_idx < nomineeList[_nomineeAddress].noArray.length, "Requested address is out of range.");
     return (nomineeList[_nomineeAddress].noArray[_idx]);
 }



 function getEvictionYesVoteCount(address _nomineeAddress)  public view returns (uint count)
 {
     return (evictionList[_nomineeAddress].yesArray.length);
 }

 function getEvictionYesVoterFromIdx(address _nomineeAddress, uint _idx)  public view returns (address voter)
 {
     require (_idx < evictionList[_nomineeAddress].yesArray.length, "Requested address is out of range.");
     return (evictionList[_nomineeAddress].yesArray[_idx]);
 }


 function getEvictionNoVoteCount(address _nomineeAddress)  public view returns (uint count)
 {
     return (evictionList[_nomineeAddress].noArray.length);
 }

 function getEvictionNoVoterFromIdx(address _nomineeAddress, uint _idx) public view returns (address voter)
 {
     require (_idx < evictionList[_nomineeAddress].noArray.length, "Requested address is out of range.");
     return (evictionList[_nomineeAddress].noArray[_idx]);
 }


 function getMoniker(address _address) public view returns (bytes32 moniker)
 {
     return (monikerList[_address]);
 }

 function getCurrentNomineeVotes(address _address) public view returns (uint yes, uint no)
 {
    if (! isNominee(_address)) {return (yes, no);}
     return (nomineeList[_address].yesVotes,nomineeList[_address].noVotes);
 }


 function getCurrentEvictionVotes(address _address) public view returns (uint yes, uint no)
 {
    if (! isEvictee(_address)) {return (yes, no);}
     return (evictionList[_address].yesVotes,evictionList[_address].noVotes);
 }

}