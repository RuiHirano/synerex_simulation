import React from 'react';
import {
    Container, connectToHarmowareVis, HarmoVisLayers, MovesLayer, LineMapLayer, MovesInput, LoadingIcon, FpsDisplay, DepotsLayer, EventInfo, MovesbaseOperation, MovesBase, BasedProps
} from 'harmoware-vis';

import Controller from '../components/controller';

import * as io from 'socket.io-client';

const MAPBOX_TOKEN = process.env.MAPBOX_ACCESS_TOKEN; //Acquire Mapbox accesstoken

class App extends Container {
    constructor(props){
	super(props);
	const { setSecPerHour, setLeading, setTrailing }= props.actions;
	setSecPerHour(3600);
	setLeading(3);
	setTrailing(3);
	const socket = io();
	this.state = {
	    moveDataVisible: true,
	    moveOptionVisible: false,
	    depotOptionVisible: false,
	    heatmapVisible: false,
	    optionChange: false,
	    popup: [0,0, ''],
		linemapData: [{
			"sourcePosition": [136.973172, 35.152476, 0],
			"targetPosition": [136.984031, 35.152476, 0],
		},
		{
			"sourcePosition": [136.973172, 35.160678, 0],
			"targetPosition": [136.984031, 35.160678, 0],
		},
		{
			"sourcePosition": [136.973172, 35.152476, 0],
			"targetPosition": [136.973172, 35.160678, 0],
		},
		{
			"sourcePosition": [136.984031, 35.152476, 0],
			"targetPosition": [136.984031, 35.160678, 0],
		},

		{
			"sourcePosition": [136.981014, 35.152476, 0],
			"targetPosition": [136.990047, 35.152476, 0],
		},
		{
			"sourcePosition": [136.981014, 35.160678, 0],
			"targetPosition": [136.990047, 35.160678, 0],
		},
		{
			"sourcePosition": [136.981014, 35.152476, 0],
			"targetPosition": [136.981014, 35.160678, 0],
		},
		{
			"sourcePosition": [136.990047, 35.152476, 0],
			"targetPosition": [136.990047, 35.160678, 0],
		},

		{
			"sourcePosition": [136.982500, 35.152476, 0],
			"targetPosition": [136.982500, 35.160678, 0],
			"color": [255, 0, 255]
		}
		]
	};

	// for receiving event info.
	socket.on('connect', ()=>{console.log("Socket.IO connected!")});
	socket.on('event', this.getEvent.bind(this));
	socket.on('disconnect', ()=>{console.log("Socket.IO disconnected!")});
	
    }

    getEvent(socketData){
		console.log("Get event!!", socketData)
	const {actions, movesbase, movedData} = this.props
	const {mtype, id,  lat, lon, angle, speed, area } = JSON.parse(socketData);
	console.log("area: ",area);
	const time = Date.now()/1000; // set time as now. (If data have time, ..)
	let hit = false;
	const movesbasedata = [...movesbase]; // why copy !?
	const movedData2 = [...movedData]; // why copy !?
	const setMovesbase = [];
	const setMovedData = [];
	
	for( let i = 0, lengthi = movesbasedata.length; i < lengthi; i+=1){
	    //	    let setMovedata = Object.assign({}, movesbasedata[i]);
	    let setMovedata = movesbasedata[i];
			let setMovedDatai = movedData2[i];
	    if(mtype === setMovedata.mtype && id === setMovedata.id){
		hit = true;
//		const {operation } = setMovedata;
		//		const arrivaltime = time;
		setMovedata.arrivaltime = time;		
		setMovedata.operation.push({
		    elapsedtime: time,
		    position:[lon, lat, 0],
		    angle,speed
		});
//		setMovedata = Object.assign({}, setMovedata, {arrivaltime, operation});
	    }
	    setMovesbase.push(setMovedata);
			setMovedData.push(setMovedDatai);
	}
	if(!hit){
	    setMovesbase.push({
				mtype, id,
				departuretime:time,
				arrivaltime: time,
				operation: [{
		    	elapsedtime:time,
		    	position:[lon, lat, 0],
		    	angle, speed
				}]
	    });
			setMovedData.push({
				sourceColor: [255, 0, 255]
	    });
	}
	actions.updateMovesBase(setMovesbase);
	//actions.updateMovedData(setMovedData);
    }

  deleteMovebase(maxKeepSecond) {
    const { actions, animatePause, movesbase, settime } = this.props
    const movesbasedata = [...movesbase];
    const setMovesbase= [];
    let dataModify = false;
    const compareTime = settime - maxKeepSecond;
    for (let i = 0, lengthi = movesbasedata.length; i < lengthi; i += 1) {
	const { departuretime:propsdeparturetime, operation:propsoperation } = movesbasedata[i];
      let departuretime = propsdeparturetime;
      let startIndex = propsoperation.length;
      for (let j = 0, lengthj = propsoperation.length; j < lengthj; j += 1) {
        if(propsoperation[j].elapsedtime > compareTime){
          startIndex = j;
          departuretime = propsoperation[j].elapsedtime;
          break;
        }
      }
      if(startIndex === 0){
        setMovesbase.push(Object.assign({}, movesbasedata[i]));
      }else
      if(startIndex < propsoperation.length){
        setMovesbase.push(Object.assign({}, movesbasedata[i], {
            operation: propsoperation.slice(startIndex), departuretime
	}));
        dataModify = true;  
      }else{
        dataModify = true;
      }
    }
    if(dataModify){
      if(!animatePause){
        actions.setAnimatePause(true);
      }
      actions.updateMovesBase(setMovesbase);
      if(!animatePause){
        actions.setAnimatePause(false);
      }
    }
  }

    getMoveDataChecked(e){
	this.setState({ moveDataVisible: e.target.checked });
    }

    getMoveOptionChecked(e){
	this.setState({ moveOptionVisible: e.target.checked });
    }

    getDepotOptionChecked(e){
	this.setState({ depotOptionVisible: e.target.checked });
    }

    getOptionChangeChecked(e){
	this.setState({ optionChange: e.target.checked });
    }
    
    
    render() {
	const props = this.props;
	const { actions, clickedObject, inputFileName, viewport, deoptsData, loading,
		routePaths, lightSettings, movesbase, movedData } = props;
//	var movedData2 = [...movedData]
//	const { movesFileName } = inputFileName;
	const optionVisible = false;
	const onHover = (el) => {
	    if (el && el.object) {
		let disptext = '';
		const objctlist = Object.entries(el.object);
		for (let i = 0, lengthi = objctlist.length; i < lengthi; i += 1) {
		    const strvalue = objctlist[i][1].toString();
		    disptext += i > 0 ? '\n' : '';
		    disptext += (`${objctlist[i][0]}: ${strvalue}`);
		}
		this.setState({ popup: [el.x, el.y, disptext] });
	    } else {
		this.setState({ popup: [0, 0, ''] });
	    }
	};
	

	return (
		<div>
		<Controller {...props}
	    deleteMovebase={this.deleteMovebase.bind(this)}
	    getMoveDataChecked ={this.getMoveDataChecked.bind(this)}
	    getMoveOptionChecked ={this.getMoveOptionChecked.bind(this)}
	    getDepotOptionChecked ={this.getDepotOptionChecked.bind(this)}
	    getOptionChangeChecked ={this.getOptionChangeChecked.bind(this)}
		/>
		<div className="harmovis_area">
		<HarmoVisLayers
	    viewport={viewport} actions={actions}
	    mapboxApiAccessToken={MAPBOX_TOKEN}
	    layers={
			this.state.moveDataVisible && movedData.length > 0 ?
			[
				new LineMapLayer( {viewport,  linemapData: this.state.linemapData }),
		    	new MovesLayer({ viewport, routePaths, movesbase, movedData,
				     clickedObject, actions, lightSettings,
				     visible: this.state.moveDataVisible,
				     optionVisible: this.state.moveOptionVisible,
				     optionChange: this.state.optionChange,
				     onHover})
			]:[
				new LineMapLayer( {viewport, linemapData: this.state.linemapData }),
			]
		}/>
		</div>
		<svg width={viewport.width} height={viewport.height} className="harmovis_overlay">
		<g fill="white" fontSize="12">
		{this.state.popup[2].length > 0 ?
		 this.state.popup[2].split('\n').map((value, index) =>
						     <text
						     x={this.state.popup[0] + 10} y={this.state.popup[1] + (index * 12)}
						     key={index.toString()}
						     >{value}</text>) : null
		}
            </g>
		</svg>
		        <LoadingIcon loading={loading} />
        <FpsDisplay />
	    </div>
	);
    }
}
export default connectToHarmowareVis(App);