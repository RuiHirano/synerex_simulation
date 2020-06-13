
import React, { useEffect } from 'react';
import { HarmoVisLayers, Container, BasedProps, BasedState, connectToHarmowareVis, MovesLayer, Movesbase, MovesbaseOperation, DepotsLayer, DepotsData, LineMapLayer } from 'harmoware-vis';
import io from "socket.io-client";
import { Controller } from '../components';


//const MAPBOX_TOKEN = process.env.REACT_APP_MAPBOX_TOKEN ? process.env.REACT_APP_MAPBOX_TOKEN : "";
const MAPBOX_TOKEN = 'pk.eyJ1IjoicnVpaGlyYW5vIiwiYSI6ImNqdmc0bXJ0dTAzZDYzem5vMmk0ejQ0engifQ.3k045idIb4JNvawjppzqZA'


class Harmoware extends Container<BasedProps & BasedState> {
    render() {
        const { actions, depotsData, viewport } = this.props;
        return (<HarmowarePage {...this.props} />)
    }
}

const socket: SocketIOClient.Socket = io();

const HarmowarePage: React.FC<BasedProps> = (props) => {
    const { actions, depotsData, viewport, movesbase, movedData, routePaths, clickedObject } = props


    const getAgents = (data: any) => {
        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        console.log("socketData length2", data.length);
        //console.log("movesbasedata length", movesbasedata.length)
        const setMovesbase: any[] = [];

        data.forEach((value: any) => {
            const { mtype, id, lat, lon, angle, speed, area } = JSON.parse(
                value
            );
            //console.log("data: ", value);

            let color = [0, 200, 0];
            if (mtype == 0) {
                // Ped
                color = [0, 200, 120];
            } else if (mtype == 1) {
                // Car
                color = [200, 0, 0];
            }
            setMovesbase.push({
                mtype,
                id,
                departuretime: time,
                arrivaltime: time,
                operation: [
                    {
                        elapsedtime: time,
                        position: [lon, lat, 0],
                        radius: 20,
                        angle,
                        speed,
                        color
                    }
                ]
            });
        });

        //console.log("lenth before", setMovesbase.length);
        actions.updateMovesBase(setMovesbase);
        //actions.updateMovedData(setMovedData);
    }

    const testAgent = async () => {
        for (let i = 0; i < 100; i++) {
            const setMovesbase: Movesbase[] = [];
            const time = Date.now() / 1000;
            let color = [0, 200, 0];
            await timeout(1000)
            for (let index = 0; index < 100; index++) {
                setMovesbase.push({
                    type: 'ped',
                    movesbaseidx: index,
                    departuretime: time,
                    arrivaltime: time,
                    operation: [
                        {
                            elapsedtime: time,
                            position: [135.4664 + index * 0.0001, 35.633253 + index * 0.0001, 0],
                            direction: 10,
                            color
                        }
                    ]
                });
            }

            actions.updateMovesBase(setMovesbase);
        }
    }

    useEffect(() => {
        socket.on("event", (data: any) => getAgents(data));
        testAgent()
        console.log(process.env);
        if (actions) {
            actions.setViewport({
                ...props.viewport,
                width: window.screen.width,
                height: window.screen.height,
                zoom: 10
            })
            actions.setSecPerHour(1000);
        }
    }, [])
    return (
        <div>
            <Controller {...props} />
            <HarmoVisLayers
                viewport={viewport} actions={actions}
                mapboxApiAccessToken={MAPBOX_TOKEN}
                layers={[
                    /*new LineLayer({
                        data: [{
                            "sourcePosition": [136.901961, 35.156615, 0],
                            "targetPosition": [136.933907, 35.144681, 0]
                        }]
                    })*/
                    new DepotsLayer({
                        depotsData,
                        iconChange: false,
                        layerRadiusScale: 20
                    }),
                    new MovesLayer({
                        routePaths,
                        movesbase,
                        movedData,
                        clickedObject,
                        actions,
                        //lightSettings,
                        layerRadiusScale: 0.1,
                        getRadius: x => 1,
                        getRouteWidth: x => 1,
                        optionCellSize: 2,
                        sizeScale: 1,
                        iconChange: false,
                        optionChange: false, // this.state.optionChange,
                        //onHover
                    })
                ]}
            />
        </div>
    );
}

async function timeout(ms: number) {
    await new Promise(resolve => setTimeout(resolve, ms));
    return
}

export default connectToHarmowareVis(Harmoware);