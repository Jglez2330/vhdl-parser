D"78"
SB"99"
12UB"X1"
12D"13"
array'LEFT
-------------------------------------------------------------------------------
-- uart
-- iMPLEMENTS A UNIVERSAL ASYNCHRONOUS RECEIVER TRANSMITTER
-------------------------------------------------------------------------------
-- CLOCK
--      iNPUT CLOCK, MUST MATCH FREQUENCY VALUE GIVEN ON CLOCK_FREQUENCY
--      GENERIC INPUT.
-- RESET
--      sYNCHRONOUS RESET.  
-- DATA_STREAM_IN
--      iNPUT DATA BUS FOR BYTES TO TRANSMIT.
-- DATA_STREAM_IN_STB
--      iNPUT STROBE TO QUALIFY THE INPUT DATA BUS.
-- DATA_STREAM_IN_ACK
--      oUTPUT ACKNOWLEDGE TO INDICATE THE uart HAS BEGUN SENDING THE BYTE
--      PROVIDED ON THE DATA_STREAM_IN PORT.
-- DATA_STREAM_OUT
--      dATA OUTPUT PORT FOR RECEIVED BYTES.
-- DATA_STREAM_OUT_STB
--      oUTPUT STROBE TO QUALIFY THE RECEIVED BYTE. wILL BE VALID FOR ONE CLOCK
--      CYCLE ONLY. 
-- TX
--      sERIAL TRANSMIT.
-- RX
--      sERIAL RECEIVE
-------------------------------------------------------------------------------
LIBRARY IEEE;
    USE IEEE.STD_LOGIC_1164.ALL;
    USE IEEE.NUMERIC_STD.ALL;
    USE IEEE.MATH_REAL.ALL;

ENTITY UART IS
    GENERIC (
        BAUD                : POSITIVE;
        CLOCK_FREQUENCY     : POSITIVE
    );
    PORT (  
        CLOCK               :   IN  STD_LOGIC;
        RESET               :   IN  STD_LOGIC;    
        DATA_STREAM_IN      :   IN  STD_LOGIC_VECTOR(7 DOWNTO 0);
        DATA_STREAM_IN_STB  :   IN  STD_LOGIC;
        DATA_STREAM_IN_ACK  :   OUT STD_LOGIC;
        DATA_STREAM_OUT     :   OUT STD_LOGIC_VECTOR(7 DOWNTO 0);
        DATA_STREAM_OUT_STB :   OUT STD_LOGIC;
        TX                  :   OUT STD_LOGIC;
        RX                  :   IN  STD_LOGIC
    );
END UART;

ARCHITECTURE RTL OF UART IS
    ---------------------------------------------------------------------------
    -- bAUD GENERATION CONSTANTS
    ---------------------------------------------------------------------------
    CONSTANT C_TX_DIV       : INTEGER := CLOCK_FREQUENCY / BAUD;
    CONSTANT C_RX_DIV       : INTEGER := CLOCK_FREQUENCY / (BAUD * 16);
    CONSTANT C_TX_DIV_WIDTH : INTEGER 
        := INTEGER(LOG2(REAL(C_TX_DIV))) + 1;   
    CONSTANT C_RX_DIV_WIDTH : INTEGER 
        := INTEGER(LOG2(REAL(C_RX_DIV))) + 1;
    ---------------------------------------------------------------------------
    -- bAUD GENERATION SIGNALS
    ---------------------------------------------------------------------------
    SIGNAL TX_BAUD_COUNTER : UNSIGNED(C_TX_DIV_WIDTH - 1 DOWNTO 0) 
        := (OTHERS => '0');   
    SIGNAL TX_BAUD_TICK : STD_LOGIC := '0';
    SIGNAL RX_BAUD_COUNTER : UNSIGNED(C_RX_DIV_WIDTH - 1 DOWNTO 0) 
        := (OTHERS => '0');   
    SIGNAL RX_BAUD_TICK : STD_LOGIC := '0';
    ---------------------------------------------------------------------------
    -- tRANSMITTER SIGNALS
    ---------------------------------------------------------------------------
    TYPE UART_TX_STATES IS ( 
        TX_SEND_START_BIT,
        TX_SEND_DATA,
        TX_SEND_STOP_BIT
    );             
    SIGNAL UART_TX_STATE : UART_TX_STATES := TX_SEND_START_BIT;
    SIGNAL UART_TX_DATA_VEC : STD_LOGIC_VECTOR(7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL UART_TX_DATA : STD_LOGIC := '1';
    SIGNAL UART_TX_COUNT : UNSIGNED(2 DOWNTO 0) := (OTHERS => '0');
    SIGNAL UART_RX_DATA_IN_ACK : STD_LOGIC := '0';
    ---------------------------------------------------------------------------
    -- rECEIVER SIGNALS
    ---------------------------------------------------------------------------
    TYPE UART_RX_STATES IS ( 
        RX_GET_START_BIT, 
        RX_GET_DATA, 
        RX_GET_STOP_BIT
    );            
    SIGNAL UART_RX_STATE : UART_RX_STATES := RX_GET_START_BIT;
    SIGNAL UART_RX_BIT : STD_LOGIC := '1';
    SIGNAL UART_RX_DATA_VEC : STD_LOGIC_VECTOR(7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL UART_RX_DATA_SR : STD_LOGIC_VECTOR(1 DOWNTO 0) := (OTHERS => '1');
    SIGNAL UART_RX_FILTER : UNSIGNED(1 DOWNTO 0) := (OTHERS => '1');
    SIGNAL UART_RX_COUNT : UNSIGNED(2 DOWNTO 0) := (OTHERS => '0');
    SIGNAL UART_RX_DATA_OUT_STB : STD_LOGIC := '0';
    SIGNAL UART_RX_BIT_SPACING : UNSIGNED (3 DOWNTO 0) := (OTHERS => '0');
    SIGNAL UART_RX_BIT_TICK : STD_LOGIC := '0';
BEGIN
    -- cONNECT io
    DATA_STREAM_IN_ACK  <= UART_RX_DATA_IN_ACK;
    DATA_STREAM_OUT     <= UART_RX_DATA_VEC;
    DATA_STREAM_OUT_STB <= UART_RX_DATA_OUT_STB;
    TX                  <= UART_TX_DATA;
    ---------------------------------------------------------------------------
    -- oversample_clock_divider
    -- GENERATE AN OVERSAMPLED TICK (BAUD * 16)
    ---------------------------------------------------------------------------
    OVERSAMPLE_CLOCK_DIVIDER : PROCESS (CLOCK)
    BEGIN
        IF RISING_EDGE (CLOCK) THEN
            IF RESET = '1' THEN
                RX_BAUD_COUNTER <= (OTHERS => '0');
                RX_BAUD_TICK <= '0';    
            ELSE
                IF RX_BAUD_COUNTER = C_RX_DIV THEN
                    RX_BAUD_COUNTER <= (OTHERS => '0');
                    RX_BAUD_TICK <= '1';
                ELSE
                    RX_BAUD_COUNTER <= RX_BAUD_COUNTER + 1;
                    RX_BAUD_TICK <= '0';
                END IF;
            END IF;
        END IF;
    END PROCESS OVERSAMPLE_CLOCK_DIVIDER;
    ---------------------------------------------------------------------------
    -- rxd_synchronise
    -- sYNCHRONISE RXD TO THE OVERSAMPLED BAUD
    ---------------------------------------------------------------------------
    RXD_SYNCHRONISE : PROCESS(CLOCK)
    BEGIN
        IF RISING_EDGE(CLOCK) THEN
            IF RESET = '1' THEN
                UART_RX_DATA_SR <= (OTHERS => '1');
            ELSE
                IF RX_BAUD_TICK = '1' THEN
                    UART_RX_DATA_SR(0) <= RX;
                    UART_RX_DATA_SR(1) <= UART_RX_DATA_SR(0);
                END IF;
            END IF;
        END IF;
    END PROCESS RXD_SYNCHRONISE;
    ---------------------------------------------------------------------------
    -- rxd_filter
    -- fILTER RXD WITH A 2 BIT COUNTER.
    ---------------------------------------------------------------------------
    RXD_FILTER : PROCESS(CLOCK)
    BEGIN
        IF RISING_EDGE(CLOCK) THEN
            IF RESET = '1' THEN
                UART_RX_FILTER <= (OTHERS => '1');
                UART_RX_BIT <= '1';
            ELSE
                IF RX_BAUD_TICK = '1' THEN
                    -- FILTER RXD.
                    IF UART_RX_DATA_SR(1) = '1' AND UART_RX_FILTER < 3 THEN
                        UART_RX_FILTER <= UART_RX_FILTER + 1;
                    ELSIF UART_RX_DATA_SR(1) = '0' AND UART_RX_FILTER > 0 THEN
                        UART_RX_FILTER <= UART_RX_FILTER - 1;
                    END IF;
                    -- SET THE RX BIT.
                    IF UART_RX_FILTER = 3 THEN
                        UART_RX_BIT <= '1';
                    ELSIF UART_RX_FILTER = 0 THEN
                        UART_RX_BIT <= '0';
                    END IF;
                END IF;
            END IF;
        END IF;
    END PROCESS RXD_FILTER;
    ---------------------------------------------------------------------------
    -- rx_bit_spacing
    ---------------------------------------------------------------------------
    RX_BIT_SPACING : PROCESS (CLOCK)
    BEGIN
        IF RISING_EDGE(CLOCK) THEN
            UART_RX_BIT_TICK <= '0';
            IF RX_BAUD_TICK = '1' THEN       
                IF UART_RX_BIT_SPACING = 15 THEN
                    UART_RX_BIT_TICK <= '1';
                    UART_RX_BIT_SPACING <= (OTHERS => '0');
                ELSE
                    UART_RX_BIT_SPACING <= UART_RX_BIT_SPACING + 1;
                END IF;
                IF UART_RX_STATE = RX_GET_START_BIT THEN
                    UART_RX_BIT_SPACING <= (OTHERS => '0');
                END IF; 
            END IF;
        END IF;
    END PROCESS RX_BIT_SPACING;
    ---------------------------------------------------------------------------
    -- uart_receive_data
    ---------------------------------------------------------------------------
    UART_RECEIVE_DATA   : PROCESS(CLOCK)
    BEGIN
        IF RISING_EDGE(CLOCK) THEN
            IF RESET = '1' THEN
                UART_RX_STATE <= RX_GET_START_BIT;
                UART_RX_DATA_VEC <= (OTHERS => '0');
                UART_RX_COUNT <= (OTHERS => '0');
                UART_RX_DATA_OUT_STB <= '0';
            ELSE
                UART_RX_DATA_OUT_STB <= '0';
                CASE UART_RX_STATE IS
                    WHEN RX_GET_START_BIT =>
                        IF RX_BAUD_TICK = '1' AND UART_RX_BIT = '0' THEN
                            UART_RX_STATE <= RX_GET_DATA;
                        END IF;
                    WHEN RX_GET_DATA =>
                        IF UART_RX_BIT_TICK = '1' THEN
                            UART_RX_DATA_VEC(UART_RX_DATA_VEC'HIGH) 
                                <= UART_RX_BIT;
                            UART_RX_DATA_VEC(
                                UART_RX_DATA_VEC'HIGH-1 DOWNTO 0
                            ) <= UART_RX_DATA_VEC(
                                UART_RX_DATA_VEC'HIGH DOWNTO 1
                            );
                            IF UART_RX_COUNT < 7 THEN
                                UART_RX_COUNT   <= UART_RX_COUNT + 1;
                            ELSE
                                UART_RX_COUNT <= (OTHERS => '0');
                                UART_RX_STATE <= RX_GET_STOP_BIT;
                            END IF;
                        END IF;
                    WHEN RX_GET_STOP_BIT =>
                        IF UART_RX_BIT_TICK = '1' THEN
                            IF UART_RX_BIT = '1' THEN
                                UART_RX_STATE <= RX_GET_START_BIT;
                                UART_RX_DATA_OUT_STB <= '1';
                            END IF;
                        END IF;                            
                    WHEN OTHERS =>
                        UART_RX_STATE <= RX_GET_START_BIT;
                END CASE;
            END IF;
        END IF;
    END PROCESS UART_RECEIVE_DATA;
    ---------------------------------------------------------------------------
    -- tx_clock_divider
    -- gENERATE BAUD TICKS AT THE REQUIRED RATE BASED ON THE INPUT CLOCK
    -- FREQUENCY AND BAUD RATE
    ---------------------------------------------------------------------------
    TX_CLOCK_DIVIDER : PROCESS (CLOCK)
    BEGIN
        IF RISING_EDGE (CLOCK) THEN
            IF RESET = '1' THEN
                TX_BAUD_COUNTER <= (OTHERS => '0');
                TX_BAUD_TICK <= '0';    
            ELSE
                IF TX_BAUD_COUNTER = C_TX_DIV THEN
                    TX_BAUD_COUNTER <= (OTHERS => '0');
                    TX_BAUD_TICK <= '1';
                ELSE
                    TX_BAUD_COUNTER <= TX_BAUD_COUNTER + 1;
                    TX_BAUD_TICK <= '0';
                END IF;
            END IF;
        END IF;
    END PROCESS TX_CLOCK_DIVIDER;
    ---------------------------------------------------------------------------
    -- uart_send_data 
    -- gET DATA FROM DATA_STREAM_IN AND SEND IT ONE BIT AT A TIME UPON EACH 
    -- BAUD TICK. sEND DATA LSB FIRST.
    -- WAIT 1 TICK, SEND START BIT (0), SEND DATA 0-7, SEND STOP BIT (1)
    ---------------------------------------------------------------------------
    UART_SEND_DATA : PROCESS(CLOCK)
    BEGIN
        IF RISING_EDGE(CLOCK) THEN
            IF RESET = '1' THEN
                UART_TX_DATA <= '1';
                UART_TX_DATA_VEC <= (OTHERS => '0');
                UART_TX_COUNT <= (OTHERS => '0');
                UART_TX_STATE <= TX_SEND_START_BIT;
                UART_RX_DATA_IN_ACK <= '0';
            ELSE
                UART_RX_DATA_IN_ACK <= '0';
                CASE UART_TX_STATE IS
                    WHEN TX_SEND_START_BIT =>
                        IF TX_BAUD_TICK = '1' AND DATA_STREAM_IN_STB = '1' THEN
                            UART_TX_DATA  <= '0';
                            UART_TX_STATE <= TX_SEND_DATA;
                            UART_TX_COUNT <= (OTHERS => '0');
                            UART_RX_DATA_IN_ACK <= '1';
                            UART_TX_DATA_VEC <= DATA_STREAM_IN;
                        END IF;
                    WHEN TX_SEND_DATA =>
                        IF TX_BAUD_TICK = '1' THEN
                            UART_TX_DATA <= UART_TX_DATA_VEC(0);
                            UART_TX_DATA_VEC(
                                UART_TX_DATA_VEC'HIGH-1 DOWNTO 0
                            ) <= UART_TX_DATA_VEC(
                                UART_TX_DATA_VEC'HIGH DOWNTO 1
                            );
                            IF UART_TX_COUNT < 7 THEN
                                UART_TX_COUNT <= UART_TX_COUNT + 1;
                            ELSE
                                UART_TX_COUNT <= (OTHERS => '0');
                                UART_TX_STATE <= TX_SEND_STOP_BIT;
                            END IF;
                        END IF;
                    WHEN TX_SEND_STOP_BIT =>
                        IF TX_BAUD_TICK = '1' THEN
                            UART_TX_DATA <= '1';
                            UART_TX_STATE <= TX_SEND_START_BIT;
                        END IF;
                    WHEN OTHERS =>
                        UART_TX_DATA <= '1';
                        UART_TX_STATE <= TX_SEND_START_BIT;
                END CASE;
            END IF;
        END IF;
    END PROCESS UART_SEND_DATA;    
END RTL;
