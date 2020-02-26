import React, { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import io from "socket.io-client";
const SequenceDiagram = require('react-sequence-diagram');

const input =
  'Andrew->China: Says Hello\n' +
  'Note over end: China thinks\\nabout it\n' +
  'China-->Andrew: How are you?\n' +
  'Andrew->>China: I am good thanks!';

const options = {
  theme: 'simple'
};

const simTypes: string[] = [
  "GET_AGENTS",
  "SET_AGENTS",
  "REGIST_PROVIDER",
  "KILL_PROVIDER",
  "DIVIDE_PROVIDER",
  "UPDATE_PROVIDERS",
  "SEND_PROVIDER_STATUS",
  "SET_PROVIDERS",
  "GET_PROVIDERS",
  "GET_CLOCK",
  "SET_CLOCK",
  "UPDATE_CLOCK",
  "START_CLOCK",
  "STOP_CLOCK",
  "FORWARD_CLOCK",
  "BACK_CLOCK"
]

interface TimeStamp {
  Seconds: number
  Nanos: number
}

class Message{
  ID: number
  Name: string
  Type: string // RD, PS, SP
  SimType: string // xxRequest, xxResponse
  TimeStamp: TimeStamp
  Targets: number[]
  SynerexAddress: string

  constructor(id: number, msgType: string, simType: string, name:string, targets: number[], synerex: string, timeStamp: TimeStamp){
    this.ID = id
    this.Type = msgType
    this.TimeStamp = timeStamp
    this.Name = name
    this.Targets = targets
    this.SynerexAddress = synerex
    this.SimType = simType
  }

  getContent() {
    return this.SimType
  }
  getInput() {
    let result = ""
    this.Targets.forEach((target: number)=>{
      result += this.ID + "->" + target + ": " + this.Type + " " + this.SimType + " " + this.TimeStamp.Seconds + "\n"

    })
    return result
  }
}

const socket: SocketIOClient.Socket = io();

/*const mockMessages: Message[] = [
  new Message(0, "RD", simTypes[0], "A", [1, 2], ":3000", 1),
  new Message(1, "PS", simTypes[1], "AB", [2], ":4000", 2),
  new Message(2, "PS", simTypes[2], "ACCC", [0, 2], ":5000", 3),
]*/

function onError(error:any) {
  console.log(error);
}

const App: React.FC = () => {

  const [store, setStore] = useState<Message[]>([]) // total message store
  const [newMessages, setNewMessages] = useState<Message[]>([]) // same time message store
  const [providers, setProviders] = useState<{[n: number]:Message}>({})
  const [inputs, setInputs] = useState<string>("")

  const addNewMessage = (mes: Message)=>{
    // newMessが0またはtimeStampが一緒であれば追加
    if(newMessages.length == 0 || mes.TimeStamp.Seconds == newMessages[newMessages.length-1].TimeStamp.Seconds){
      // setNewMessage
      setNewMessages(prevMes =>{
        let mess = prevMes
        mess.push(mes)
        // sort
        mess.sort(function(a, b) {
          if (a.TimeStamp.Nanos > b.TimeStamp.Nanos) {
            return 1;
          } else {
            return -1;
          }
        })
        console.log("mes: ", mess)
        return mess
      })
    }else if(mes.TimeStamp.Seconds > newMessages[newMessages.length-1].TimeStamp.Seconds){
      const newInputs = getNewInputs()
      setInputs(prevInputs => {
        let cpInputs = prevInputs // copy
        cpInputs += newInputs
        return cpInputs
      })

      // 新しいMesをいれてNewMessをClear
      setNewMessages([mes])

    }
  }

  useEffect(()=>{

    socket.on("connect", () => {
      console.log("Socket.IO connected!");
    });

    socket.on("event", (data: string) => {
      const data2 = data.split(',\"arg\":\"')
      const mainInfoJson = data2[0] + "}"
      const senderInfoJson = data2[1].substring(0, data2[1].length-2)
      
      
      const mainInfo: any = JSON.parse(mainInfoJson)
      const senderInfo: any = JSON.parse(senderInfoJson)
      let mes: Message
      const ts: TimeStamp = {
        Seconds: senderInfo.ts.seconds,
        Nanos: senderInfo.ts.nanos
      }
      if(mainInfo.msgType === "RD"){
        const simDemand = senderInfo.ArgOneof.SimDemand
        mes = new Message(simDemand.pid, mainInfo.msgType, simTypes[simDemand.type], simDemand.sender_info.name, simDemand.targets, simDemand.sender_info.synerex_address, ts)
      }else{
        const simSupply = senderInfo.ArgOneof.SimSupply
        mes = new Message(simSupply.pid, mainInfo.msgType, simTypes[simSupply.type], simSupply.sender_info.name, simSupply.targets, simSupply.sender_info.synerex_address, ts)
      }

      addNewMessage(mes)
    });

    socket.on("disconnect", () => {
        console.log("Socket.IO disconnected!");
    });

  }, [])

  const getNewInputs = ()=>{
    // create new Inputs
    let newInputs = "" 
    newInputs += "Note over EndFlag: Finish\n"
    newMessages.forEach((mes: Message, index: number)=>{
      newInputs += mes.getInput()
    })
    return newInputs
  }
  console.log("render")

  return (
    <div className="App">
      <text>Synerex Sequence Diagram</text>
      <SequenceDiagram input={inputs + getNewInputs()} options={options} onError={onError} />
    </div>
  );
}

export default App;
