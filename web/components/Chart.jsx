const { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ReferenceLine } = Recharts;

const SqueezeChart = () => {
    const data = [
        { time: '1', close: 10, value: 0, sqzOn: true },
        { time: '2', close: 12, value: 0.2, sqzOn: true },
        { time: '3', close: 15, value: 0.5, sqzOn: false },
        { time: '4', close: 14, value: 0.3, sqzOn: false },
        { time: '5', close: 13, value: -0.1, sqzOn: true },
        { time: '6', close: 11, value: -0.3, sqzOn: true },
        { time: '7', close: 10, value: -0.4, sqzOn: false },
        { time: '8', close: 9, value: -0.5, sqzOn: false },
        { time: '9', close: 10, value: 0.1, sqzOn: true },
        { time: '10', close: 11, value: 0.4, sqzOn: true }
    ];

    return (
        <div className="w-full p-4">
            <div className="mb-4 text-xl font-bold">Squeeze Indicator</div>
            <LineChart width={800} height={400} data={data} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis />
                <Tooltip />
                <Legend />
                <ReferenceLine y={0} stroke="#666" />
                <Line
                    type="monotone"
                    dataKey="close"
                    stroke="#8884d8"
                    name="Price"
                />
                <Line
                    type="monotone"
                    dataKey="value"
                    stroke="#82ca9d"
                    name="Squeeze Value"
                    dot={({ payload }) => (
                        <circle
                            cx={0}
                            cy={0}
                            r={4}
                            fill={payload.sqzOn ? "#ffd700" : "#82ca9d"}
                        />
                    )}
                />
            </LineChart>
        </div>
    );
};